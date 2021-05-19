group "default" {
    targets = ["mutating-ac", "validating-ac"]
}

variable "TAG" {
    default = "latest"
}

target "mutating-ac" {
    dockerfile = "mutating-ac/Dockerfile"
    context = "./"
    tags = ["ghcr.io/IIJ-Global-Solutions-Vietnam/namespace-admission-controller/mutating-ac:${TAG}"]
}

target "validating-ac" {
    dockerfile = "validating-ac/Dockerfile"
    context = "./"
    tags = ["ghcr.io/IIJ-Global-Solutions-Vietnam/namespace-admission-controller/validating-ac:${TAG}"]
}
