load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@bazel_gazelle//:def.bzl", "gazelle")

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
        "//posts",
        "@com_github_gorilla_handlers//:go_default_library",
        "@com_github_gorilla_mux//:go_default_library",
        "@com_github_gorilla_sessions//:go_default_library",
    ],
)

go_binary(
    name = "quillpen",
    embed = [":quillpen_lib"],
    visibility = ["//visibility:public"],
)
