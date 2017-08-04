# Part 5: Add providers

## What's needed to be done.
 - disk provider
 - aws provider

## Code - Interface.

In part4, we upload the file on disk. But to be `cloud ready`, we need to implement
an other type of storage.  
In this exemple, this is the AWS' S3 storage.  
To achieve that, we will create an interface. This allows us to abstract the way we want to
handle the project's upload part.  

```go
    type Provider interface {
        Get(filename string) (io.ReadCloser, error)
        Put(filename string, image multipart.File) error
    }
```
The disk provider and the aws provider must implement this interface.

This way, we can add all the storage provider witch are on the market.
An other one, could be the blob storage of the Google Cloud.

