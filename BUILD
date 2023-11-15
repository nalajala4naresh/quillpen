load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@bazel_gazelle//:def.bzl", "gazelle")
load("@rules_oci//oci:defs.bzl", "oci_image","oci_push")
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")

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
        "//pkg/accounts",
        "//pkg/chat",
        "//pkg/posts",
        "@com_github_gorilla_handlers//:handlers",
        "@com_github_gorilla_mux//:mux",
        "@com_github_gorilla_sessions//:sessions",
    ],
)

go_binary(
    name = "quillpen",
    embed = [":quillpen_lib"],
    visibility = ["//visibility:public"],
)



pkg_tar(
    name = "tar",
    srcs = [":quillpen"],
)


oci_image(
    name = "quillpen_image",
    base = "@distroless_base",
    tars = [":tar"],
    entrypoint = ["/quillpen"],
)

oci_push(
    name = "push_quillpen",
    image = ":quillpen_image",
    repository = "index.docker.io/nalajalanaresh/quillpen",
    remote_tags = ["latest"]

)

alias(
    name = "go",
    actual = "@go_sdk//:bin/go",
    visibility = ["//visibility:public"],
)