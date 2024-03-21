# JENKINS-BUILD-INFO

A small GitHub action for obtaining the image name with tag from jenkins build info.

## Example usage
```yaml
    - name: Get Image name
      uses: ahton89/jenkins-build-info@v0.0.1
      with:
        image_name: example
```

## Inputs
```yaml
    github_token:
      description: 'Github Token'
      required: false
    image_name:
      description: 'image name from jenkins'
      required: true
```

## Outputs
```yaml
  jenkins-image-name:
    description: 'Image name with tag from jenkins build info'
```

No warranty is given, use at your own risk. 🤷‍
