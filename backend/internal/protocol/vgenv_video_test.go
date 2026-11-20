package protocol

import (
	"context"
	"reflect"
	"strings"
	"testing"
)

func vgenvVideoAdapter(t *testing.T) Adapter {
	t.Helper()
	return officialPackageAdapter(t, "vgenv-video.yingce-plugin", "vgenv-video")
}

func TestVgenvVideoCreateTextToVideo(t *testing.T) {
	adapter := vgenvVideoAdapter(t)
	request := GenerationRequest{
		Model:       "base:fast",
		Prompt:      "动漫少女在樱花树下微笑，阳光明媚，微风吹拂",
		Duration:    5,
		AspectRatio: "9:16",
		Resolution:  "768p",
		Extra: map[string]any{
			"idempotencyKey": "vgenv-test-key-123",
		},
	}
	spec, err := adapter.BuildCreate(context.Background(), RequestContext{
		BaseURL: "https://api.vgenv.com",
		Request: request,
	})
	if err != nil {
		t.Fatalf("BuildCreate failed: %v", err)
	}
	if spec.Method != "POST" || spec.Path != "/v3/videos" || spec.Headers["Idempotency-Key"] != "vgenv-test-key-123" {
		t.Fatalf("spec metadata mismatch: %#v", spec)
	}
	body := manifestTestBody(t, spec)
	want := map[string]any{
		"model":        "base",
		"tier":         "fast",
		"prompt":       "动漫少女在樱花树下微笑，阳光明媚，微风吹拂",
		"duration":     float64(5),
		"resolution":   "768p",
		"aspect_ratio": "9:16",
	}
	if !reflect.DeepEqual(body, want) {
		t.Fatalf("body = %#v, want %#v", body, want)
	}
}

func TestVgenvVideoCreateFirstAndLastFrame(t *testing.T) {
	adapter := vgenvVideoAdapter(t)
	request := GenerationRequest{
		Model:       "cinematic:standard",
		Prompt:      "镜头从 @图片1 的特写平滑推移过渡到 @图片2 的全景画面",
		Duration:    6,
		AspectRatio: "16:9",
		Resolution:  "768p",
		Images: []MediaReference{
			{ID: "node-1", Kind: "image", Role: "first_frame", URL: "https://example.com/first.png", Order: 0},
			{ID: "node-2", Kind: "image", Role: "last_frame", URL: "https://example.com/last.png", Order: 1},
		},
		Extra: map[string]any{
			"idempotencyKey": "key-frames-456",
		},
	}
	spec, err := adapter.BuildCreate(context.Background(), RequestContext{
		BaseURL: "https://api.vgenv.com",
		Request: request,
	})
	if err != nil {
		t.Fatalf("BuildCreate failed: %v", err)
	}
	body := manifestTestBody(t, spec)
	expectedPrompt := "镜头从 [First frame] 的特写平滑推移过渡到 [Last frame] 的全景画面"
	if body["prompt"] != expectedPrompt {
		t.Fatalf("prompt = %q, want %q", body["prompt"], expectedPrompt)
	}
	frameImages, ok := body["frame_images"].([]any)
	if !ok || len(frameImages) != 2 {
		t.Fatalf("frame_images = %#v, want 2 items", body["frame_images"])
	}
	first := frameImages[0].(map[string]any)
	if first["frame_type"] != "first_frame" || first["image_url"].(map[string]any)["url"] != "https://example.com/first.png" {
		t.Fatalf("first frame mismatch: %#v", first)
	}
	last := frameImages[1].(map[string]any)
	if last["frame_type"] != "last_frame" || last["image_url"].(map[string]any)["url"] != "https://example.com/last.png" {
		t.Fatalf("last frame mismatch: %#v", last)
	}
	if _, exists := body["reference_images"]; exists {
		t.Fatalf("reference_images should not exist in frame_images mode: %#v", body["reference_images"])
	}
}

func TestVgenvVideoCreateReferenceImages(t *testing.T) {
	adapter := vgenvVideoAdapter(t)
	request := GenerationRequest{
		Model:       "base:fast",
		Prompt:      "保持主角 @图片1 的发型和服饰，场景参考 @图片2 的夜市街道氛围",
		Duration:    5,
		AspectRatio: "9:16",
		Resolution:  "768p",
		Images: []MediaReference{
			{ID: "char-1", Kind: "image", Role: "reference_image", URL: "https://example.com/char.png", Order: 0},
			{ID: "scene-1", Kind: "image", Role: "reference_image", URL: "https://example.com/scene.png", Order: 1},
		},
		Extra: map[string]any{
			"idempotencyKey": "key-ref-789",
		},
	}
	spec, err := adapter.BuildCreate(context.Background(), RequestContext{
		BaseURL: "https://api.vgenv.com",
		Request: request,
	})
	if err != nil {
		t.Fatalf("BuildCreate failed: %v", err)
	}
	body := manifestTestBody(t, spec)
	expectedPrompt := "保持主角 <Picture 1> 的发型和服饰，场景参考 <Picture 2> 的夜市街道氛围"
	if body["prompt"] != expectedPrompt {
		t.Fatalf("prompt = %q, want %q", body["prompt"], expectedPrompt)
	}
	refImages, ok := body["reference_images"].([]any)
	if !ok || len(refImages) != 2 {
		t.Fatalf("reference_images = %#v, want 2 items", body["reference_images"])
	}
	ref1 := refImages[0].(map[string]any)
	if ref1["image_url"].(map[string]any)["url"] != "https://example.com/char.png" {
		t.Fatalf("reference image 1 mismatch: %#v", ref1)
	}
	ref2 := refImages[1].(map[string]any)
	if ref2["image_url"].(map[string]any)["url"] != "https://example.com/scene.png" {
		t.Fatalf("reference image 2 mismatch: %#v", ref2)
	}
	if _, exists := body["frame_images"]; exists {
		t.Fatalf("frame_images should not exist when only references provided: %#v", body["frame_images"])
	}
}

func TestVgenvVideoCreateFirstFrameAndReferenceImages(t *testing.T) {
	adapter := vgenvVideoAdapter(t)
	// 通过 videoStartFrameNodeId 显式指定 startFrameNodeId 为 node-start，另外一张为参考图
	request := GenerationRequest{
		Model:       "cinematic:fast",
		Prompt:      "以 @图片1 为起始画面，画面中的人物动作和道具风格参考 @图片2",
		Duration:    5,
		AspectRatio: "16:9",
		Resolution:  "768p",
		Images: []MediaReference{
			{ID: "node-start", Kind: "image", URL: "https://example.com/start.png", Order: 0},
			{ID: "node-ref", Kind: "image", URL: "https://example.com/ref.png", Order: 1},
		},
		Extra: map[string]any{
			"idempotencyKey":        "key-first-and-ref",
			"videoStartFrameNodeId": "node-start",
		},
	}
	spec, err := adapter.BuildCreate(context.Background(), RequestContext{
		BaseURL: "https://api.vgenv.com",
		Request: request,
	})
	if err != nil {
		t.Fatalf("BuildCreate failed: %v", err)
	}
	body := manifestTestBody(t, spec)
	expectedPrompt := "以 [First frame] 为起始画面，画面中的人物动作和道具风格参考 <Picture 1>"
	if body["prompt"] != expectedPrompt {
		t.Fatalf("prompt = %q, want %q", body["prompt"], expectedPrompt)
	}
	frameImages, ok := body["frame_images"].([]any)
	if !ok || len(frameImages) != 1 {
		t.Fatalf("frame_images = %#v, want 1 item", body["frame_images"])
	}
	first := frameImages[0].(map[string]any)
	if first["frame_type"] != "first_frame" || first["image_url"].(map[string]any)["url"] != "https://example.com/start.png" {
		t.Fatalf("first frame mismatch: %#v", first)
	}
	refImages, ok := body["reference_images"].([]any)
	if !ok || len(refImages) != 1 {
		t.Fatalf("reference_images = %#v, want 1 item", body["reference_images"])
	}
	ref1 := refImages[0].(map[string]any)
	if ref1["image_url"].(map[string]any)["url"] != "https://example.com/ref.png" {
		t.Fatalf("ref image mismatch: %#v", ref1)
	}
}

func TestVgenvVideoPollAndResponse(t *testing.T) {
	adapter := vgenvVideoAdapter(t)
	pollSpec, err := adapter.BuildPoll(context.Background(), PollContext{
		BaseURL: "https://api.vgenv.com",
		TaskID:  "vid_task_123456",
	})
	if err != nil {
		t.Fatalf("BuildPoll failed: %v", err)
	}
	if pollSpec.Method != "GET" || pollSpec.Path != "/v3/videos/vid_task_123456" {
		t.Fatalf("pollSpec = %#v", pollSpec)
	}

	rawResponse := `{
		"id": "vid_task_123456",
		"status": "succeeded",
		"output": {
			"url": "https://cdn.vgenv.com/videos/output-123456.mp4"
		},
		"price": {
			"currency": "USD",
			"amount_micros": 50000
		}
	}`
	state, err := adapter.ParsePoll(context.Background(), PollContext{
		TaskID: "vid_task_123456",
	}, []byte(rawResponse))
	if err != nil {
		t.Fatalf("ParsePoll failed: %v", err)
	}
	if state.Status != StatusSucceeded || state.TaskID != "vid_task_123456" {
		t.Fatalf("state metadata mismatch: %#v", state)
	}
	if state.Result == nil || len(state.Result.Videos) != 1 || state.Result.Videos[0].URL != "https://cdn.vgenv.com/videos/output-123456.mp4" {
		t.Fatalf("state.Result.Videos = %#v", state.Result)
	}
	if !state.Result.Videos[0].Ephemeral {
		t.Fatalf("video should be ephemeral")
	}
}

func TestVgenvVideoValidations(t *testing.T) {
	adapter := vgenvVideoAdapter(t)
	// 尾帧不能与参考图同时使用
	reqInvalid := GenerationRequest{
		Model:       "base:fast",
		Prompt:      "test prompt",
		Duration:    5,
		AspectRatio: "16:9",
		Resolution:  "768p",
		Images: []MediaReference{
			{ID: "img-1", Kind: "image", Role: "last_frame", URL: "https://example.com/last.png"},
			{ID: "img-2", Kind: "image", Role: "reference_image", URL: "https://example.com/ref.png"},
		},
	}
	_, err := adapter.BuildCreate(context.Background(), RequestContext{
		BaseURL: "https://api.vgenv.com",
		Request: reqInvalid,
	})
	if err == nil || !strings.Contains(err.Error(), "尾帧不能与参考图同时使用") {
		t.Fatalf("expected error about last_frame and reference conflict, got: %v", err)
	}

	// Origin 模型不支持首帧+参考图
	reqOrigin := GenerationRequest{
		Model:       "origin:fast",
		Prompt:      "test origin",
		Duration:    5,
		AspectRatio: "16:9",
		Resolution:  "768p",
		Images: []MediaReference{
			{ID: "img-1", Kind: "image", Role: "first_frame", URL: "https://example.com/first.png"},
			{ID: "img-2", Kind: "image", Role: "reference_image", URL: "https://example.com/ref.png"},
		},
	}
	_, err = adapter.BuildCreate(context.Background(), RequestContext{
		BaseURL: "https://api.vgenv.com",
		Request: reqOrigin,
	})
	if err == nil || !strings.Contains(err.Error(), "Origin 不支持“首帧 + 参考图”") {
		t.Fatalf("expected error about Origin not supporting first frame + reference, got: %v", err)
	}

	// 提示词小于 5 个字符
	reqShortPrompt := GenerationRequest{
		Model:       "base:fast",
		Prompt:      "跳起来",
		Duration:    5,
		AspectRatio: "16:9",
		Resolution:  "768p",
	}
	_, err = adapter.BuildCreate(context.Background(), RequestContext{
		BaseURL: "https://api.vgenv.com",
		Request: reqShortPrompt,
	})
	if err == nil || !strings.Contains(err.Error(), "提示词至少需要 5 个字符") {
		t.Fatalf("expected error about prompt length, got: %v", err)
	}
}

func TestVgenvVideoSingleImageReferenceOperation(t *testing.T) {
	adapter := vgenvVideoAdapter(t)
	request := GenerationRequest{
		Model:       "base:fast",
		Prompt:      "主角保持 @图片1 的动漫形象，在海滩漫步",
		Duration:    5,
		AspectRatio: "9:16",
		Resolution:  "768p",
		Operation:   "reference_to_video",
		Images: []MediaReference{
			{ID: "char-1", Kind: "image", URL: "https://example.com/char.png", Order: 0},
		},
		Extra: map[string]any{
			"idempotencyKey": "key-single-ref",
		},
	}
	spec, err := adapter.BuildCreate(context.Background(), RequestContext{
		BaseURL: "https://api.vgenv.com",
		Request: request,
	})
	if err != nil {
		t.Fatalf("BuildCreate failed: %v", err)
	}
	body := manifestTestBody(t, spec)
	expectedPrompt := "主角保持 <Picture 1> 的动漫形象，在海滩漫步"
	if body["prompt"] != expectedPrompt {
		t.Fatalf("prompt = %q, want %q", body["prompt"], expectedPrompt)
	}
	refImages, ok := body["reference_images"].([]any)
	if !ok || len(refImages) != 1 {
		t.Fatalf("reference_images = %#v, want 1 item", body["reference_images"])
	}
	if _, exists := body["frame_images"]; exists {
		t.Fatalf("frame_images should not exist when operation is reference_to_video")
	}
}
