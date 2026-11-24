# vgenv Video 接口字段

## 协议身份

- 插件 ID：`vgenv-video`。
- Provider ID：`vgenv-video`。
- 能力：`video`。
- 默认 Base URL：`https://api.vgenv.com`。
- 鉴权驱动：`bearer`。
- 创建：`POST /v3/videos`。
- 查询：`GET /v3/videos/{{taskId}}`。
- 取消：`DELETE /v3/videos/{{taskId}}`。

## 配置字段

| 字段 | 类型 | 必填 | 含义 |
| --- | --- | --- | --- |
| `apiKey` | secret | 是 | API Key |

## 统一字段映射

| 统一字段 | 类型 | 必填 | 上游映射 | 说明 |
| --- | --- | --- | --- | --- |
| `model` | string | 是 | `model` | 视频模型 ID。 |
| `prompt` | string | 是 | `prompt/content/input` | 视频提示词。 |
| `images` | media[] | 否 | `first/last/reference image` | 显式 role 图片输入。 |
| `videos` | media[] | 否 | `reference video` | 参考视频。 |
| `audios` | media[] | 否 | `reference audio/voice` | 参考音频或音色。 |
| `duration` | integer | 否 | `duration/seconds` | 时长秒数。 |
| `aspectRatio` | string | 否 | `ratio/aspect_ratio/size` | 画幅比例或尺寸。 |
| `resolution` | string | 否 | `resolution` | 分辨率档位。 |
| `generateAudio` | boolean | 否 | `generate_audio` | 是否生成音频。 |
| `watermark` | boolean | 否 | `watermark` | 水印开关。 |
| `providerOptions` | object | 否 | `provider-specific fields` | 插件命名空间内的厂商扩展字段。 |

## 上游请求模板逐字段清单

下表由插件请求模板生成，覆盖 body、query、headers 和 multipart 文件声明中的每个字段。

| 上游位置 | 值或转换表达式 |
| --- | --- |
| `create.method` | `"POST"` |
| `create.path` | `"/v3/videos"` |
| `create.contentType` | `"application/json"` |
| `create.body.model` | `{"$at":[{"$split":[{"$ref":"request.model"},":"]},0]}` |
| `create.body.tier` | `{"$coalesce":[{"$at":[{"$split":[{"$ref":"request.model"},":"]},1]},"fast"]}` |
| `create.body.prompt` | `{"$ref":"request.prompt"}` |
| `create.body.duration` | `{"$omitEmpty":{"$if":{"condition":{"$gt":[{"$ref":"request.duration"},0]},"then":{"$ref":"request.duration"},"else":null}}}` |
| `create.body.resolution` | `{"$omitEmpty":{"$lower":{"$trim":{"$ref":"request.resolution"}}}}` |
| `create.body.aspect_ratio` | `{"$omitEmpty":{"$trim":{"$ref":"request.aspectRatio"}}}` |
| `create.body.seed` | `{"$omitEmpty":{"$ref":"request.providerOptions.vgenv-video.seed"}}` |
| `create.body.prompt_optimization` | `{"$omitEmpty":{"$if":{"condition":{"$ref":"request.providerOptions.vgenv-video.prompt_optimization"},"then":{"mode":{"$ref":"request.providerOptions.vgenv-video.prompt_optimization"}},"else":null}}}` |
| `create.body.combat_enhancement` | `{"$omitEmpty":{"$if":{"condition":{"$ref":"request.providerOptions.vgenv-video.combat_enhancement"},"then":true,"else":null}}}` |
| `create.body.motion_speed` | `{"$omitEmpty":{"$if":{"condition":{"$ref":"request.providerOptions.vgenv-video.combat_enhancement"},"then":{"$ref":"request.providerOptions.vgenv-video.motion_speed"},"else":null}}}` |
| `create.body.upscale` | `{"$omitEmpty":{"$if":{"condition":{"$eq":[{"$ref":"request.providerOptions.vgenv-video.upscale"},"2x"]},"then":"2x","else":null}}}` |
| `create.body.frame_images` | `{"$omitEmpty":{"$concatArrays":[{"$map":{"from":{"$filter":{"from":{"$sortByOrder":{"$ref":"request.images"}},"as":"media","where":{"$in":[{"$ref":"media.role"},["first_frame"]]}}},"as":"media","in":{"type":"image_url","image_url":{"url":{"$ref":"media.value"}},"frame_type":"first_frame"}}},{"$map":{"from":{"$filter":{"from":{"$sortByOrder":{"$ref":"request.images"}},"as":"media","where":{"$in":[{"$ref":"media.role"},["last_frame"]]}}},"as":"media","in":{"type":"image_url","image_url":{"url":{"$ref":"media.value"}},"frame_type":"last_frame"}}}]}}` |
| `create.body.reference_images` | `{"$omitEmpty":{"$map":{"from":{"$filter":{"from":{"$sortByOrder":{"$ref":"request.images"}},"as":"media","where":{"$in":[{"$ref":"media.role"},["reference_image","subject_reference","style_reference",""]]}}},"as":"media","in":{"type":"image_url","image_url":{"url":{"$ref":"media.value"}}}}}}` |
| `create.headers.Idempotency-Key` | `{"$ref":"request.extra.idempotencyKey"}` |
| `poll.method` | `"GET"` |
| `poll.path` | `"/v3/videos/{{taskId}}"` |
| `cancel.method` | `"DELETE"` |
| `cancel.path` | `"/v3/videos/{{taskId}}"` |

## Provider 扩展键

- `providerOptions.vgenv-video.combat_enhancement`
- `providerOptions.vgenv-video.motion_speed`
- `providerOptions.vgenv-video.prompt_optimization`
- `providerOptions.vgenv-video.seed`
- `providerOptions.vgenv-video.upscale`

动态模型或工作流允许使用文档声明的完整 `parameters/input/extra_body` 对象；该对象是协议本身的开放 schema，不会被宿主裁剪。

## 响应映射逐字段清单

| 映射位置 | 上游路径或转换表达式 |
| --- | --- |
| `response.taskId` | `{"$coalesce":[{"$ref":"response.id"},{"$ref":"taskId"}]}` |
| `response.status` | `{"$coalesce":[{"$ref":"response.status"},"queued"]}` |
| `response.message` | `{"$coalesce":[{"$ref":"response.error.message"},{"$ref":"response.message"}]}` |
| `response.stage` | `{"$if":{"condition":{"$eq":[{"$ref":"response.status"},"preparing"]},"then":"准备素材与审核提示词","else":null}}` |
| `response.queueAhead` | `{"$ref":"response.queue_position"}` |
| `response.videos` | `{"$ref":"response.output.url"}` |
| `response.usage` | `{"$omitEmpty":{"priceCurrency":{"$ref":"response.price.currency"},"priceAmountMicros":{"$ref":"response.price.amount_micros"}}}` |
| `response.errorPaths[0]` | `"error.code"` |
| `response.resultEphemeral` | `true` |

## 响应与错误

插件把上游 task/status/text/media/usage 映射为统一结果。临时媒体 URL 标记为 ephemeral，由宿主立即下载持久化。HTTP 错误、业务 code 和 error object 保持失败语义，不包装成成功。

## 兼容边界

模型 ID 写成 <model>:<tier>。创建返回 HTTP 202，轮询 GET /v3/videos/{id}，状态 preparing/queued/running 为进行中，succeeded 取 output.url（签名地址 24 小时有效，宿主立即下载转存），failed/cancelled 为终态；只有 preparing/queued 可取消。Idempotency-Key 由宿主任务提交键提供，同一键重试不会重复冻结余额。图片只接受 JPEG/PNG/WebP 公网 HTTPS 地址或 Data URL，单张不超过 5 MiB。

<!-- YINGCE_MANIFEST_CONTRACT_START -->
## Manifest 完整接口定义

以下 JSON 与插件包内实际 `manifest.json` 逐字段一致，覆盖插件身份、权限、配置、鉴权、参数、校验、创建、Agent、查询、取消、结果下载、响应和 Agent 响应映射。`documentation` 字段的值就是当前完整文档；为避免文档在自身内部无限递归，JSON 中仅用等义占位文本表示正文。

```json
{
  "apiVersion": "yingce.plugin/v2",
  "id": "vgenv-video",
  "name": "vgenv Video",
  "version": "2.1.0",
  "author": "vgenv / 影策",
  "description": "vgenv V3 异步视频接口：文生视频、首帧、尾帧、首尾帧、全能参考（最多 9 张图）与首帧 + 参考图。",
  "documentation": "<当前插件的完整 documentation，由 README.md 与 docs/interface.md 拼接而成；为避免 JSON 递归，此处不重复展开正文。>",
  "permissions": [
    "generation.run",
    "media.read"
  ],
  "configuration": {
    "fields": [
      {
        "name": "apiKey",
        "type": "secret",
        "label": "API Key",
        "required": true
      }
    ]
  },
  "contributes": {
    "providers": [
      {
        "id": "vgenv-video",
        "label": "vgenv Video",
        "capabilities": [
          "video"
        ],
        "scopes": [
          "admin.system-channel",
          "user.custom-channel",
          "canvas",
          "creation",
          "agent"
        ],
        "baseUrl": "https://api.vgenv.com",
        "requiresPublicMediaUrls": false,
        "detachedPolling": true,
        "idempotentCreate": true,
        "auth": {
          "type": "bearer",
          "field": "apiKey"
        },
        "videoReferences": {
          "layout": "reference-images",
          "token": "<Picture {n}>",
          "dramaTailFrame": "reference-image"
        },
        "parameters": [
          {
            "name": "model",
            "type": "string",
            "required": true,
            "mapping": "model",
            "description": "视频模型 ID。"
          },
          {
            "name": "prompt",
            "type": "string",
            "required": true,
            "mapping": "prompt/content/input",
            "description": "视频提示词。"
          },
          {
            "name": "images",
            "type": "media[]",
            "required": false,
            "mapping": "first/last/reference image",
            "description": "显式 role 图片输入。"
          },
          {
            "name": "videos",
            "type": "media[]",
            "required": false,
            "mapping": "reference video",
            "description": "参考视频。"
          },
          {
            "name": "audios",
            "type": "media[]",
            "required": false,
            "mapping": "reference audio/voice",
            "description": "参考音频或音色。"
          },
          {
            "name": "duration",
            "type": "integer",
            "required": false,
            "mapping": "duration/seconds",
            "description": "时长秒数。"
          },
          {
            "name": "aspectRatio",
            "type": "string",
            "required": false,
            "mapping": "ratio/aspect_ratio/size",
            "description": "画幅比例或尺寸。"
          },
          {
            "name": "resolution",
            "type": "string",
            "required": false,
            "mapping": "resolution",
            "description": "分辨率档位。"
          },
          {
            "name": "generateAudio",
            "type": "boolean",
            "required": false,
            "mapping": "generate_audio",
            "description": "是否生成音频。"
          },
          {
            "name": "watermark",
            "type": "boolean",
            "required": false,
            "mapping": "watermark",
            "description": "水印开关。"
          },
          {
            "name": "providerOptions",
            "type": "object",
            "required": false,
            "mapping": "provider-specific fields",
            "description": "插件命名空间内的厂商扩展字段。"
          }
        ],
        "validations": [
          {
            "assert": {
              "$gte": [
                {
                  "$len": {
                    "$trim": {
                      "$ref": "request.prompt"
                    }
                  }
                },
                5
              ]
            },
            "message": "vgenv 提示词至少需要 5 个字符，请补充动作、场景或角色描述"
          },
          {
            "assert": {
              "$in": [
                {
                  "$at": [
                    {
                      "$split": [
                        {
                          "$ref": "request.model"
                        },
                        ":"
                      ]
                    },
                    0
                  ]
                },
                [
                  "base",
                  "cinematic",
                  "origin",
                  "minimax/h3",
                  "minimax/h3-human",
                  "minimax/h3-native"
                ]
              ]
            },
            "message": "vgenv 模型 ID 必须写成 <model>:<tier>，如 base:fast、cinematic:fast、origin:fast 或 minimax/h3:fast"
          },
          {
            "assert": {
              "$in": [
                {
                  "$coalesce": [
                    {
                      "$at": [
                        {
                          "$split": [
                            {
                              "$ref": "request.model"
                            },
                            ":"
                          ]
                        },
                        1
                      ]
                    },
                    "fast"
                  ]
                },
                [
                  "fast",
                  "standard",
                  "enhanced",
                  "quality"
                ]
              ]
            },
            "message": "vgenv 档位必须是 fast、standard、enhanced 或 quality"
          },
          {
            "assert": {
              "$or": [
                {
                  "$eq": [
                    {
                      "$ref": "request.duration"
                    },
                    0
                  ]
                },
                {
                  "$and": [
                    {
                      "$gte": [
                        {
                          "$ref": "request.duration"
                        },
                        2
                      ]
                    },
                    {
                      "$lte": [
                        {
                          "$ref": "request.duration"
                        },
                        15
                      ]
                    }
                  ]
                }
              ]
            },
            "message": "vgenv 视频时长应为 2–15 秒整数（480p 为 4–15 秒）"
          },
          {
            "assert": {
              "$in": [
                {
                  "$lower": {
                    "$trim": {
                      "$coalesce": [
                        {
                          "$ref": "request.resolution"
                        },
                        "480p"
                      ]
                    }
                  }
                },
                [
                  "480p",
                  "768p"
                ]
              ]
            },
            "message": "vgenv 分辨率只支持 480p 或 768p"
          },
          {
            "assert": {
              "$in": [
                {
                  "$coalesce": [
                    {
                      "$trim": {
                        "$ref": "request.aspectRatio"
                      }
                    },
                    "9:16"
                  ]
                },
                [
                  "9:16",
                  "2:3",
                  "3:4",
                  "1:1",
                  "4:3",
                  "3:2",
                  "16:9",
                  "21:9"
                ]
              ]
            },
            "message": "vgenv 不支持该画幅比例"
          },
          {
            "assert": {
              "$lte": [
                {
                  "$len": {
                    "$filter": {
                      "from": {
                        "$sortByOrder": {
                          "$ref": "request.images"
                        }
                      },
                      "as": "media",
                      "where": {
                        "$in": [
                          {
                            "$ref": "media.role"
                          },
                          [
                            "reference_image",
                            "subject_reference",
                            "style_reference",
                            ""
                          ]
                        ]
                      }
                    }
                  }
                },
                9
              ]
            },
            "message": "vgenv 最多支持 9 张参考图"
          },
          {
            "assert": {
              "$lte": [
                {
                  "$len": {
                    "$filter": {
                      "from": {
                        "$sortByOrder": {
                          "$ref": "request.images"
                        }
                      },
                      "as": "media",
                      "where": {
                        "$in": [
                          {
                            "$ref": "media.role"
                          },
                          [
                            "first_frame"
                          ]
                        ]
                      }
                    }
                  }
                },
                1
              ]
            },
            "message": "vgenv 首帧最多 1 张"
          },
          {
            "assert": {
              "$lte": [
                {
                  "$len": {
                    "$filter": {
                      "from": {
                        "$sortByOrder": {
                          "$ref": "request.images"
                        }
                      },
                      "as": "media",
                      "where": {
                        "$in": [
                          {
                            "$ref": "media.role"
                          },
                          [
                            "last_frame"
                          ]
                        ]
                      }
                    }
                  }
                },
                1
              ]
            },
            "message": "vgenv 尾帧最多 1 张"
          },
          {
            "assert": {
              "$eq": [
                {
                  "$len": {
                    "$ref": "request.videos"
                  }
                },
                0
              ]
            },
            "message": "vgenv 参考视频需要先上传为素材，当前插件暂不支持，请移除参考视频"
          },
          {
            "assert": {
              "$eq": [
                {
                  "$len": {
                    "$ref": "request.audios"
                  }
                },
                0
              ]
            },
            "message": "vgenv 参考音频需要先上传为素材，当前插件暂不支持，请移除参考音频"
          },
          {
            "assert": {
              "$not": {
                "$and": [
                  {
                    "$gt": [
                      {
                        "$len": {
                          "$filter": {
                            "from": {
                              "$sortByOrder": {
                                "$ref": "request.images"
                              }
                            },
                            "as": "media",
                            "where": {
                              "$in": [
                                {
                                  "$ref": "media.role"
                                },
                                [
                                  "last_frame"
                                ]
                              ]
                            }
                          }
                        }
                      },
                      0
                    ]
                  },
                  {
                    "$gt": [
                      {
                        "$len": {
                          "$filter": {
                            "from": {
                              "$sortByOrder": {
                                "$ref": "request.images"
                              }
                            },
                            "as": "media",
                            "where": {
                              "$in": [
                                {
                                  "$ref": "media.role"
                                },
                                [
                                  "reference_image",
                                  "subject_reference",
                                  "style_reference",
                                  ""
                                ]
                              ]
                            }
                          }
                        }
                      },
                      0
                    ]
                  }
                ]
              }
            },
            "message": "vgenv 尾帧不能与参考图同时使用；请改用“首帧 + 参考图”，或去掉参考图"
          },
          {
            "assert": {
              "$or": [
                {
                  "$ne": [
                    {
                      "$at": [
                        {
                          "$split": [
                            {
                              "$ref": "request.model"
                            },
                            ":"
                          ]
                        },
                        0
                      ]
                    },
                    "origin"
                  ]
                },
                {
                  "$eq": [
                    {
                      "$len": {
                        "$filter": {
                          "from": {
                            "$sortByOrder": {
                              "$ref": "request.images"
                            }
                          },
                          "as": "media",
                          "where": {
                            "$in": [
                              {
                                "$ref": "media.role"
                              },
                              [
                                "reference_image",
                                "subject_reference",
                                "style_reference",
                                ""
                              ]
                            ]
                          }
                        }
                      }
                    },
                    0
                  ]
                },
                {
                  "$eq": [
                    {
                      "$len": {
                        "$filter": {
                          "from": {
                            "$sortByOrder": {
                              "$ref": "request.images"
                            }
                          },
                          "as": "media",
                          "where": {
                            "$in": [
                              {
                                "$ref": "media.role"
                              },
                              [
                                "first_frame"
                              ]
                            ]
                          }
                        }
                      }
                    },
                    0
                  ]
                }
              ]
            },
            "message": "vgenv Origin 不支持“首帧 + 参考图”，请改用 Base 或 Cinematic"
          }
        ],
        "create": {
          "method": "POST",
          "path": "/v3/videos",
          "contentType": "application/json",
          "body": {
            "model": {
              "$at": [
                {
                  "$split": [
                    {
                      "$ref": "request.model"
                    },
                    ":"
                  ]
                },
                0
              ]
            },
            "tier": {
              "$coalesce": [
                {
                  "$at": [
                    {
                      "$split": [
                        {
                          "$ref": "request.model"
                        },
                        ":"
                      ]
                    },
                    1
                  ]
                },
                "fast"
              ]
            },
            "prompt": {
              "$ref": "request.prompt"
            },
            "duration": {
              "$omitEmpty": {
                "$if": {
                  "condition": {
                    "$gt": [
                      {
                        "$ref": "request.duration"
                      },
                      0
                    ]
                  },
                  "then": {
                    "$ref": "request.duration"
                  },
                  "else": null
                }
              }
            },
            "resolution": {
              "$omitEmpty": {
                "$lower": {
                  "$trim": {
                    "$ref": "request.resolution"
                  }
                }
              }
            },
            "aspect_ratio": {
              "$omitEmpty": {
                "$trim": {
                  "$ref": "request.aspectRatio"
                }
              }
            },
            "seed": {
              "$omitEmpty": {
                "$ref": "request.providerOptions.vgenv-video.seed"
              }
            },
            "prompt_optimization": {
              "$omitEmpty": {
                "$if": {
                  "condition": {
                    "$ref": "request.providerOptions.vgenv-video.prompt_optimization"
                  },
                  "then": {
                    "mode": {
                      "$ref": "request.providerOptions.vgenv-video.prompt_optimization"
                    }
                  },
                  "else": null
                }
              }
            },
            "combat_enhancement": {
              "$omitEmpty": {
                "$if": {
                  "condition": {
                    "$ref": "request.providerOptions.vgenv-video.combat_enhancement"
                  },
                  "then": true,
                  "else": null
                }
              }
            },
            "motion_speed": {
              "$omitEmpty": {
                "$if": {
                  "condition": {
                    "$ref": "request.providerOptions.vgenv-video.combat_enhancement"
                  },
                  "then": {
                    "$ref": "request.providerOptions.vgenv-video.motion_speed"
                  },
                  "else": null
                }
              }
            },
            "upscale": {
              "$omitEmpty": {
                "$if": {
                  "condition": {
                    "$eq": [
                      {
                        "$ref": "request.providerOptions.vgenv-video.upscale"
                      },
                      "2x"
                    ]
                  },
                  "then": "2x",
                  "else": null
                }
              }
            },
            "frame_images": {
              "$omitEmpty": {
                "$concatArrays": [
                  {
                    "$map": {
                      "from": {
                        "$filter": {
                          "from": {
                            "$sortByOrder": {
                              "$ref": "request.images"
                            }
                          },
                          "as": "media",
                          "where": {
                            "$in": [
                              {
                                "$ref": "media.role"
                              },
                              [
                                "first_frame"
                              ]
                            ]
                          }
                        }
                      },
                      "as": "media",
                      "in": {
                        "type": "image_url",
                        "image_url": {
                          "url": {
                            "$ref": "media.value"
                          }
                        },
                        "frame_type": "first_frame"
                      }
                    }
                  },
                  {
                    "$map": {
                      "from": {
                        "$filter": {
                          "from": {
                            "$sortByOrder": {
                              "$ref": "request.images"
                            }
                          },
                          "as": "media",
                          "where": {
                            "$in": [
                              {
                                "$ref": "media.role"
                              },
                              [
                                "last_frame"
                              ]
                            ]
                          }
                        }
                      },
                      "as": "media",
                      "in": {
                        "type": "image_url",
                        "image_url": {
                          "url": {
                            "$ref": "media.value"
                          }
                        },
                        "frame_type": "last_frame"
                      }
                    }
                  }
                ]
              }
            },
            "reference_images": {
              "$omitEmpty": {
                "$map": {
                  "from": {
                    "$filter": {
                      "from": {
                        "$sortByOrder": {
                          "$ref": "request.images"
                        }
                      },
                      "as": "media",
                      "where": {
                        "$in": [
                          {
                            "$ref": "media.role"
                          },
                          [
                            "reference_image",
                            "subject_reference",
                            "style_reference",
                            ""
                          ]
                        ]
                      }
                    }
                  },
                  "as": "media",
                  "in": {
                    "type": "image_url",
                    "image_url": {
                      "url": {
                        "$ref": "media.value"
                      }
                    },
                    "url": {
                      "$ref": "media.value"
                    }
                  }
                }
              }
            }
          },
          "headers": {
            "Idempotency-Key": {
              "$ref": "request.extra.idempotencyKey"
            }
          }
        },
        "poll": {
          "method": "GET",
          "path": "/v3/videos/{{taskId}}"
        },
        "cancel": {
          "method": "DELETE",
          "path": "/v3/videos/{{taskId}}"
        },
        "response": {
          "taskId": {
            "$coalesce": [
              {
                "$ref": "response.id"
              },
              {
                "$ref": "taskId"
              }
            ]
          },
          "status": {
            "$coalesce": [
              {
                "$ref": "response.status"
              },
              "queued"
            ]
          },
          "message": {
            "$coalesce": [
              {
                "$ref": "response.error.message"
              },
              {
                "$ref": "response.message"
              }
            ]
          },
          "stage": {
            "$if": {
              "condition": {
                "$eq": [
                  {
                    "$ref": "response.status"
                  },
                  "preparing"
                ]
              },
              "then": "准备素材与审核提示词",
              "else": null
            }
          },
          "queueAhead": {
            "$ref": "response.queue_position"
          },
          "videos": {
            "$ref": "response.output.url"
          },
          "usage": {
            "$omitEmpty": {
              "priceCurrency": {
                "$ref": "response.price.currency"
              },
              "priceAmountMicros": {
                "$ref": "response.price.amount_micros"
              }
            }
          },
          "errorPaths": [
            "error.code"
          ],
          "resultEphemeral": true
        }
      }
    ]
  }
}
```
<!-- YINGCE_MANIFEST_CONTRACT_END -->
