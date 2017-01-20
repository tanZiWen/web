# Prosnav OAuth2.0 RESTApi

### Request authorize code 
### Common parameters
parameter | value
-- | --
host | http://locahost:9090

#### URL: 
/web/v1/oauth/authorize

#### Method:
GET

#### Parameters:

| Parameter     |  type  | required |
|---------------|:------:|----------|
| response_type | string | true     |
| client_id     | string | true     |
| redirect_uri  | string | true     |
| state         | string | false    |

#### Example:

<pre>
URL:  localhost:9090/oauth/v1/authorize?response_type=code&client_id=123&redirect_uri=http%3A%2F%2Fhuaban.com%2F   
Redirect page: https://www.baidu.com/?code=61E21Mu5SkuHRRYmou9VvQ&state=
</pre>

### Request access_token

#### URL: 
/web/v1/oauth/token

#### Parameters:

| Parameter     |  type  | required |
|---------------|:------:|----------|
| code          | string | true     |
| grant_type    | string | true     |
| client_id     | string | true     |
| client_secret | string | true     |
| redirect_uri  | string | false    |
    
#### Method

GET, POST

#### Example:

<pre>
URL:  localhost:9090/oauth/v1/token?grant_type=refresh_token
Method: POST
Parameters: 
    grant_type=authorization_code&
    code=QemJvk0gTyaO6hja1Q_yZQ&
    client_id=123&
    client_secret=123
Return:
    {
        "access_token": "IxCk6POdSZC9GROmXkqv8w",
        "expires_in": 3600,
        "refresh_token": "5sVpnweaSV2jTGgVjKQPxg",
        "token_type": "Bearer"
    }

</pre>

###Request user profile

#### URL: 
/oauth/v1/userInfo

#### Parameters:

| Parameter     |  type  | required |
|---------------|:------:|----------|
| code          | string | true     |
    
#### Method

GET, POST

#### Example:

<pre>
URL:  localhost:9090/oauth/v1/userInfo?code=m4_y2qHCTHKbiu9TKwtc-g
Method: GET
Return:
    {
        "UserId": "client"
    }
</pre>

###Refresh access token

#### URL: 
/oauth/v1/token?grant_type=refresh_token

#### Parameters:

| Parameter     |  type  | required |
|---------------|:------:|----------|
| grant_type    | string | true     |
| refresh_token | string | true     |
| client_id     | string | true     |
| client_secret | string | true     |
    
#### Method

GET, POST

#### Example:

<pre>
URL:  http://localhost:9090/oauth/v1/token?grant_type=refresh_token&refresh_token=rFG4v6jYSRavXXurEC1Zyg&client_id=client&client_secret=client
grant_type=refresh_token
Method: GET
Return:
    {
        "access_token": "Y2rXYeUTQOOiCVxDPqORtg",
        "expires_in": 3600,
        "refresh_token": "Z4mLAOiKQNqoVkw2byZX5w",
        "token_type": "Bearer"
    }
</pre>

###Request user list

#### URL: 
/oauth/v1/users

#### Parameters:
Type: json,

| Parameter | type   | required |
|-----------|:----- :|:---------|
| code      | string | true     |
| pageIndex | int    | true     |
| pageCount | int    | true     |
| userInfo  | string | false    |
| realName  | string | false    |
| area      | string | false    |
    
#### Method

POST

#### Example:

<pre>
URL:  http://192.168.1.28:9090/oauth/v1/users?code=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0YW55dWFuIiwiY2lkIjoiMTIzIiwiZXhwIjoxNDY0ODMzODA1fQ.Oh68SYwHkHqm0eWQQmENcLYY9mp5EVxSsykwv0LVDyo
Method: POST
Content-Type: application/json
Parameters: 
    {
         "pageIndex": 1,
         "pageCount": 10,
         "userInfo": "张",
    }
Return:
    {
        "count":29,
        "users":[{"_id":31,"area":"SH","position":"rm","realName":"张勇"},
            {"_id":49,"area":"SH","position":"","realName":"张羽"},
            {"_id":92,"area":"SH","position":"","realName":"张元元"},
            {"_id":166,"area":"SH","position":"","realName":"机构合作-张福钢"},
            {"_id":179,"area":"SH","position":"","realName":"渠道-张福钢"},
            {"_id":200,"area":"SH","position":"","realName":"张霞雯"},
            {"_id":220,"area":"SH","position":"","realName":"张艺琼"},
            {"_id":234,"area":"SH","position":"rm","realName":"张海霞"},
            {"_id":293,"position":"consultant","realName":"张译尹"},
            {"_id":303,"area":"SH","position":"","realName":"张琦"}]
    }
</pre>



