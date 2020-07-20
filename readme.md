# Traefik Plugin GeoIP2

> plugin in development
> also see https://github.com/negasus/traefik-plugin-ip2location

Add a request header with Country ISO code detected by IP with MaxMind GeoIP database 

## Configuration

### `Filename (string) required`

Path to the database file

### `FromHeader (string) (default: empty)`

By default, plugin obtain an IP address from a request.RemoteAddr.
If you want to get IP address from the HTTP header, you could define this option.

### `CountryHeader (string) (default: 'X-Country')`

HTTP Header name for place the Country code

## Other

If an any error occurred, Country header will be empty. 



