package registry

const defaultLifecyclePolicy = `
{
  "rules": [
    {
      "action": {
        "type": "expire"
      },
      "selection": {
        "countType": "imageCountMoreThan",
        "countNumber": 100,
        "tagStatus": "tagged",
        "tagPrefixList": [
          "master-",
          "release-"
        ]
      },
      "description": "Keep 100 primary images",
      "rulePriority": 1
    },
    {
      "action": {
        "type": "expire"
      },
      "selection": {
        "countType": "sinceImagePushed",
        "countUnit": "days",
        "countNumber": 7,
        "tagStatus": "untagged"
      },
      "description": "Remove untagged images older than a week",
      "rulePriority": 2
    },
    {
      "action": {
        "type": "expire"
      },
      "selection": {
        "countType": "imageCountMoreThan",
        "countNumber": 200,
        "tagStatus": "any"
      },
      "description": "No more than 200 images",
      "rulePriority": 10
    }
  ]
}`
