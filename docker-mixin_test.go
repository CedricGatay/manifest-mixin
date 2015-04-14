package main

import (
    "testing"
)

func Test_invalid_argument_count_error_no_arg(t *testing.T){
    t.Error("Not yet tested")
}

func Test_read_remote_file_not_found(t *testing.T){
    t.Error("Not yet tested")
}

func Test_cant_open_Dockerfile(t *testing.T){
    _, err := readRemoteDockerFile("<NOO>")
    if err == nil {
        t.Error("Should not be able to open file")
    }
}

func Test_comment_FROM_lines(t *testing.T){
    fileContent := `Blablabla
FROM Blabla
Tralala`
    cleaned := cleanupDockerFile("comment_FROM_lines", fileContent)
    expected := `##### MIXIN BEGIN (comment_FROM_lines)
Blablabla
# FROM Blabla
Tralala
##### MIXIN END 
`
    if cleaned != expected{
        t.Errorf("%s is not equal to %s", cleaned, expected)
    }
}

func Test_comment_MAINTAINER_lines(t *testing.T){
    fileContent := `Blablabla
MAINTAINER <c.gatay@code-troopers.com>
Tralala`
    cleaned := cleanupDockerFile("comment_MAINTAINER_lines", fileContent)
    expected := `##### MIXIN BEGIN (comment_MAINTAINER_lines)
Blablabla
# MAINTAINER <c.gatay@code-troopers.com>
Tralala
##### MIXIN END 
`
    if cleaned != expected{
        t.Errorf("%s is not equal to %s", cleaned, expected)
    }
}
