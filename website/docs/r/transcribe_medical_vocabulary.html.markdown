---
subcategory: "Transcribe"
layout: "aws"
page_title: "AWS: aws_transcribe_medical_vocabulary"
description: |-
  Resource for managing an AWS Transcribe MedicalVocabulary.
---

# Resource: aws_transcribe_medical_vocabulary

Resource for managing an AWS Transcribe MedicalVocabulary.

## Example Usage

### Basic Usage

```terraform
resource "aws_s3_bucket" "example" {
  bucket        = "example-medical-vocab-123"
  force_destroy = true
}

resource "aws_s3_object" "object" {
  bucket = aws_s3_bucket.example.id
  key    = "transcribe/test1.txt"
  source = "test.txt"
}

resource "aws_transcribe_medical_vocabulary" "example" {
  vocabulary_name     = "example"
  language_code       = "en-US"
  vocabulary_file_uri = "s3://${aws_s3_bucket.example.id}/${aws_s3_object.object.key}"

  tags = {
    tag1 = "value1"
    tag2 = "value3"
  }

  depends_on = [
    aws_s3_object.object
  ]
}
```

## Argument Reference

The following arguments are required:

* `language_code` - (Required) The language code you selected for your medical vocabulary. US English (en-US) is the only language supported with Amazon Transcribe Medical.
* `vocabulary_file_uri` - (Required) The Amazon S3 location (URI) of the text file that contains your custom medical vocabulary.
* `vocabulary_name` - (Required) The name of the Medical Vocabulary.

The following arguments are optional:

* `tags` - (Optional) A map of tags to assign to the MedicalVocabulary. If configured with a provider `default_tags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Name of the MedicalVocabulary.
* `arn` - ARN of the MedicalVocabulary.
* `download_uri` - Generated download URI.

## Timeouts

Configuration options:

* `create` - (Default `30m`)
* `update` - (Default `30m`)
* `delete` - (Default `30m`)

## Import

Transcribe MedicalVocabulary can be imported using the `vocabulary_name`, e.g.,

```
$ terraform import aws_transcribe_medical_vocabulary.example example-name
```
