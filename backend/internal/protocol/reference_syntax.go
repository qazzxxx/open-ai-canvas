// 声明式协议的“素材角色与提示词引用约定”：上游对首帧/尾帧/参考图的区分，以及提示词里引用某张图的写法。
//
// 画布提示词用 @图片N 引用第 N 张连接的图片（编号按连线顺序，首帧与尾帧同样计入）。
// 宿主按清单声明把它改写成上游认识的令牌，避免每个插件各自猜编号。

package protocol

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// ManifestReferenceSyntax 由插件声明，宿主在 BuildCreate 之前应用。
type ManifestReferenceSyntax struct {
	// Image 是参考图的提示词写法，{n} 为参考图序号（从 1 起，只数参考图，不含首尾帧）。
	Image string `json:"image,omitempty"`
	// FirstFrame / LastFrame 是首帧、尾帧的提示词写法；为空时保留原始 @图片N。
	FirstFrame string `json:"firstFrame,omitempty"`
	LastFrame  string `json:"lastFrame,omitempty"`
	// HonorFrameNodes 表示画布显式选择的首帧/尾帧节点（request.extra.videoStartFrameNodeId、
	// videoEndFrameNodeId）优先于宿主的参考模式推断。画布在全能参考模式下会把首帧降级为参考图，
	// 开启后首帧仍按首帧发送。
	HonorFrameNodes bool `json:"honorFrameNodes,omitempty"`
	// SingleImageIsFirstFrame 表示没有任何首尾帧时，唯一一张图作为首帧（图生视频）。
	SingleImageIsFirstFrame bool `json:"singleImageIsFirstFrame,omitempty"`
}

const (
	roleFirstFrame = "first_frame"
	roleLastFrame  = "last_frame"
	roleReference  = "reference_image"
)

// 宿主在 request.extra 中传递的首尾帧节点 ID，与 canvas 元数据字段同名。
const (
	extraStartFrameNodeID = "videoStartFrameNodeId"
	extraEndFrameNodeID   = "videoEndFrameNodeId"
)

var canvasImageMention = regexp.MustCompile(`@图片(\d+)`)

func validateReferenceSyntax(syntax *ManifestReferenceSyntax) error {
	if syntax == nil {
		return nil
	}
	for name, value := range map[string]string{"image": syntax.Image, "firstFrame": syntax.FirstFrame, "lastFrame": syntax.LastFrame} {
		if len(value) > 64 || strings.ContainsAny(value, "\r\n") {
			return fmt.Errorf("%s must be a single line of at most 64 bytes", name)
		}
	}
	if syntax.Image != "" && !strings.Contains(syntax.Image, "{n}") {
		return fmt.Errorf("image template must contain {n}")
	}
	return nil
}

func isReferenceRole(role string) bool {
	switch role {
	case roleReference, "subject_reference", "style_reference":
		return true
	}
	return false
}

// applyReferenceSyntax 返回一份新请求：每张图都有明确角色，提示词里的 @图片N 已按上游写法改写。
// 不修改调用方的请求，图片顺序保持不变。
func applyReferenceSyntax(syntax *ManifestReferenceSyntax, request GenerationRequest) (GenerationRequest, error) {
	if syntax == nil || len(request.Images) == 0 {
		return request, nil
	}
	images := append([]MediaReference(nil), request.Images...)
	// 编号按连线顺序（Order）；Order 相同时保持传入顺序。
	positions := make([]int, len(images))
	for index := range positions {
		positions[index] = index
	}
	sort.SliceStable(positions, func(i, j int) bool { return images[positions[i]].Order < images[positions[j]].Order })

	if syntax.HonorFrameNodes {
		startID := stringExtra(request, extraStartFrameNodeID)
		endID := stringExtra(request, extraEndFrameNodeID)
		for index := range images {
			switch id := images[index].ID; {
			case id != "" && id == startID:
				images[index].Role = roleFirstFrame
			case id != "" && id == endID:
				images[index].Role = roleLastFrame
			}
		}
	}
	if syntax.SingleImageIsFirstFrame && len(images) == 1 && !hasFrameRole(images) && request.Operation != "reference_to_video" {
		images[0].Role = roleFirstFrame
	}
	for index := range images {
		if role := strings.TrimSpace(images[index].Role); role != roleFirstFrame && role != roleLastFrame && !isReferenceRole(role) {
			images[index].Role = roleReference
		}
	}

	// 画布编号按原始连接位置（Order+1）；调用方没有填写不重复的 Order 时退回到传入位置。
	distinct := make(map[int]bool, len(images))
	for _, image := range images {
		distinct[image.Order] = true
	}
	labelOf := func(index int) int {
		if len(distinct) == len(images) {
			return images[index].Order + 1
		}
		return index + 1
	}
	tokens := make(map[int]string, len(images))
	referenceNumber := 0
	for _, index := range positions {
		role := strings.TrimSpace(images[index].Role)
		switch {
		case role == roleFirstFrame:
			tokens[labelOf(index)] = syntax.FirstFrame
		case role == roleLastFrame:
			tokens[labelOf(index)] = syntax.LastFrame
		case isReferenceRole(role):
			referenceNumber++
			if syntax.Image != "" {
				tokens[labelOf(index)] = strings.ReplaceAll(syntax.Image, "{n}", strconv.Itoa(referenceNumber))
			}
		}
	}
	request.Images = images
	request.Prompt = canvasImageMention.ReplaceAllStringFunc(request.Prompt, func(mention string) string {
		number, err := strconv.Atoi(canvasImageMention.FindStringSubmatch(mention)[1])
		if err != nil || tokens[number] == "" {
			return mention
		}
		return tokens[number]
	})
	return request, nil
}

func hasFrameRole(images []MediaReference) bool {
	for _, image := range images {
		if role := strings.TrimSpace(image.Role); role == roleFirstFrame || role == roleLastFrame {
			return true
		}
	}
	return false
}

func stringExtra(request GenerationRequest, key string) string {
	value, _ := request.Extra[key].(string)
	return strings.TrimSpace(value)
}
