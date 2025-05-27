# URL-Shortener

This is a url-shortener implemented in golang. It uses badger db as its key value store. It uses crc32 hash package in order to create hash for the urls.  

## Usage

Make a POST request to localhost:8080 (default port is 8080) with body
``` 
{
    "Url": "www.github.com"
}
```

Then it will return the hash alongside the address e.g. localhost:8080/12345678 where 12345678 is the hash  
Visiting localhost:8080/3697404839 will redirect us to the url provided.  


You can also visit the path /view to view all hashses, but it will write to stdout not ResponseWriter  
