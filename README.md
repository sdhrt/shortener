# URL-Shortener

This is a url-shortener implemented in golang. It uses sqlite to store data. It uses crc32 hash package in order to create hash for the urls.  

### Packages used
- modernc.org/sqlite

## Usage

Make a POST request to localhost:8080 (default port is 8080) with body
``` 
{
    "Url": "www.github.com"
}
```

Then it will return the hash alongside the address e.g. localhost:8080/12345678 where 12345678 is the hash  
Visiting localhost:8080/12345678 will redirect us to the url provided.  

## TODO
- [ ] distributed, use sharding
- [ ] Add caching service
- [ ] Load balancing
