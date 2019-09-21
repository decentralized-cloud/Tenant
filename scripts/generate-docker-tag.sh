#!/usr/bin/env bash

set -e

# making sure we change the directory to where the current script resides
current_directory=$(dirname "$0")
cd "$current_directory"

build_source_branch="$1"
build_id="$2"

if [[ "$build_source_branch" == refs/pull* ]]; then
    # building repository name triggered by Pull Request
    build_source_branch="${build_source_branch//refs\/pull\//}"
    build_source_branch="${build_source_branch//\/merge/}"
    build_source_branch="${build_source_branch//\//-}"
    tag_to_return=merge-$build_source_branch-latest
    tag_id_to_return=merge-$build_source_branch-$build_id
elif [ "$build_source_branch" = "refs/heads/master" ]; then
    # building repository name triggered by a build for the master branch
    tag_to_return="latest"
    tag_id_to_return=$build_id
else
    # building repository name triggered by a build for a working branch
    build_source_branch="${build_source_branch//refs\/heads\//}"
    build_source_branch="${build_source_branch//\//-}"
    tag_to_return=$build_source_branch-latest
    tag_id_to_return=$build_source_branch-$build_id
fi

azure_command_to_set_variable="##vso[task.setvariable variable=GENERATED_DOCKER_TAG]"
azure_command_with_value="$azure_command_to_set_variable$tag_to_return"
echo "$azure_command_with_value"

azure_command_to_set_variable="##vso[task.setvariable variable=GENERATED_DOCKER_TAG_ID]"
azure_command_with_value="$azure_command_to_set_variable$tag_id_to_return"
echo "$azure_command_with_value"
