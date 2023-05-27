load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@bazel_gazelle//:def.bzl", "gazelle")
load("@io_bazel_rules_docker//go:image.bzl", "go_image")

gazelle(
    name = "gazelle",
    prefix = "github.com/quillpen",
)

gazelle(
    name = "gazelle-update-repos",
    args = [
        "-from_file=go.mod",
        "-to_macro=deps.bzl%go_dependencies",
        "-prune",
    ],
    command = "update-repos",
)

go_library(
    name = "quillpen_lib",
    srcs = ["main.go"],
    importpath = "github.com/quillpen",
    visibility = ["//visibility:private"],
    deps = [
        "//accounts",
        "//chat",
        "//posts",
        "@com_github_gorilla_handlers//:handlers",
        "@com_github_gorilla_mux//:mux",
        "@com_github_gorilla_sessions//:sessions",
    ],
)

go_binary(
    name = "quillpen",
    data = [":html_templates"],
    embed = [":quillpen_lib"],
    visibility = ["//visibility:public"],
)

filegroup(
    name = "html_templates",
    srcs = glob(
        ["**/*.html"],
    ),
)

go_image(
    name = "quill_image",
    binary = ":quillpen",
    visibility = ["//visibility:public"],
)

alias(
    name = "go",
    actual = "@go_sdk//:bin/go",
    visibility = ["//visibility:public"],
)
