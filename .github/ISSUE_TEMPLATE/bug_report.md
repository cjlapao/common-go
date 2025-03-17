---
name: Bug report
description: Create a bug report
title: "[BUG] "
labels:

- bug
- triage
assignees:
- cjlapao
body:
- type: markdown
    attributes:
      value: |
        Thanks for taking the time to fill out this bug report!
- type: checkboxes
    attributes:
      label: Is there an existing issue for this?
      description: Search to see if an issue already exists for the bug you encountered.
      options:
  - label: I have searched the existing issues
        required: true
- type: textarea
    attributes:
      label: Current Behavior
      description: A concise description of what you're experiencing.
      required: true
- type: textarea
    attributes:
      label: Expected Behavior
      description: A concise description of what you expected to happen.
    validations:
      required: true
- type: textarea
    attributes:
      label: common-go version
      description: |
        common-go version where you observed this issue
      placeholder: |
          vX.Y.Z
      render: Markdown
    validations:
      required: true
- type: textarea
    attributes:
      label: Steps To Reproduce
      description: |
        Steps to reproduce the issue.
        To speed up the triaging of your request, reproduce the issue running
      placeholder: |
        1. In this environment...
        2. With this config...
        3. Run '...'
        4. See error...
    validations:
      required: true
- type: textarea
    attributes:
      label: Anything else?
      description: |
        Links? References? Anything that will give us more context about the issue you are encountering!

        Tip: You can attach images or log files by clicking this area to highlight it and then dragging files in.
    validations:
      required: false
- type: checkboxes
    id: validation
    attributes:
      label: Validation
      options:
          - label: Yes, I've included all of the above information (Version, settings, logs, etc.)
            required: true
