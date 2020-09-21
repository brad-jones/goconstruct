# Terraform File Example

This shows how this library might be used to represent the following terraform constructs:

- <https://registry.terraform.io/providers/hashicorp/local/latest/docs/resources/file>
- <https://registry.terraform.io/providers/hashicorp/local/latest/docs/data-sources/file>

## Expected Output

```json
{
    "id": "/",
    "type": "Root",
    "props": null,
    "children": [
        {
            "id": "/MyFileData",
            "type": "Local/FileData",
            "props": {
                "Filename": "./go.mod"
            },
            "children": []
        },
        {
            "id": "/MyFile",
            "type": "Local/File",
            "props": {
                "SensitiveContent": "",
                "ContentBase64": "",
                "Filename": "./foo.txt",
                "FilePermission": "",
                "DirectoryPermission": "",
                "Source": "",
                "Content": "/MyFileData.Content"
            },
            "children": []
        }
    ]
}
```
