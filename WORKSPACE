load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "f4a9314518ca6acfa16cc4ab43b0b8ce1e4ea64b81c38d8a3772883f153346b8",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.50.1/rules_go-v0.50.1.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.50.1/rules_go-v0.50.1.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "d76bf7a60fd8b050444090dfa2837a4eaf9829e1165618ee35dceca5cbdf58d5",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.37.0/bazel-gazelle-v0.37.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.37.0/bazel-gazelle-v0.37.0.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_download_sdk", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

############################################################
# Define your own dependencies here using go_repository.
# Else, dependencies declared by rules_go/gazelle will be used.
# The first declaration of an external repository "wins".
############################################################

go_download_sdk(
    name = "go_sdk",
    goarch = "arm64",
    goos = "darwin",
    sdks = {
        # NOTE: In most cases the whole sdks attribute is not needed.
        # There are 2 "common" reasons you might want it:
        #
        # 1. You need to use an modified GO SDK, or an unsupported version
        #    (for example, a beta or release candidate)
        #
        # 2. You want to avoid the dependency on the index file for the
        #    SHA-256 checksums. In this case, You can get the expected
        #    filenames and checksums from https://go.dev/dl/
        "linux_amd64": ("go1.20.11.linux-amd64.tar.gz", "ef79a11aa095a08772d2a69e4f152f897c4e96ee297b0dc20264b7dec2961abe"),
        "darwin_arm64": ("go1.20.11.darwin-arm64.tar.gz", "9c29373ae278dc02c1d8a08465eea8490bb09bea4ba88fce79b01fe75f84961c"),
        "linux_arm64" : ("go1.20.11.linux-arm64.tar.gz","7908a49c6ce9d48af9b5ba76ccaa0769da45d8b635259a01065b3739acef4ada"),
    },
    version = "1.20.11",
)

load("//:deps.bzl", "go_dependencies")

# gazelle:repository_macro deps.bzl%go_dependencies
go_dependencies()

go_rules_dependencies()

go_register_toolchains()

gazelle_dependencies(go_sdk = "go_sdk")


#rules_pkg



http_archive(
    name = "rules_pkg",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_pkg/releases/download/1.0.1/rules_pkg-1.0.1.tar.gz",
        "https://github.com/bazelbuild/rules_pkg/releases/download/1.0.1/rules_pkg-1.0.1.tar.gz",
    ],
    sha256 = "d20c951960ed77cb7b341c2a59488534e494d5ad1d30c4818c736d57772a9fef",
)
load("@rules_pkg//:deps.bzl", "rules_pkg_dependencies")
rules_pkg_dependencies()

# rules_oci


http_archive(
    name = "rules_oci",
    sha256 = "d007e6c96eb62c88397b68f329e4ca56e0cfe31204a2c54b0cb17819f89f83c8",
    strip_prefix = "rules_oci-2.0.0",
    url = "https://github.com/bazel-contrib/rules_oci/releases/download/v2.0.0/rules_oci-v2.0.0.tar.gz",
)

load("@rules_oci//oci:dependencies.bzl", "rules_oci_dependencies")

rules_oci_dependencies()

load("@rules_oci//oci:repositories.bzl", "LATEST_CRANE_VERSION", "LATEST_ZOT_VERSION", "oci_register_toolchains")

oci_register_toolchains(
    name = "oci",
    crane_version = LATEST_CRANE_VERSION,
    # Uncommenting the zot toolchain will cause it to be used instead of crane for some tasks.
    # Note that it does not support docker-format images.
    # zot_version = LATEST_ZOT_VERSION,
)

# You can pull your base images using oci_pull like this:
load("@rules_oci//oci:pull.bzl", "oci_pull")

oci_pull(
    name = "distroless_base",
    digest = "sha256:c925d12234f8d3fbef2256012b168004d4c47a82c4f06afcfd06fd208732fbe0",
    image = "gcr.io/distroless/base",
    platforms = [
        "linux/amd64",
        "linux/arm64",
        "darwin/arm64"
    ],
)
