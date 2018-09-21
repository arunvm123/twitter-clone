# twitter-clone

## How to run?

```bash
docker-compose up
```

## API Reference

* Signup (Method-POST)
    ```
    http://localhost:8080/signup
    ```
    Request body should contain the data as json.
    
    ```
    {
        "name": "John Doe",
        "user_name": "doe_j",
        "password": "jd7845"
    }
    ```

* Login (Method-POST)
    ```
    http://localhost:8080/signup
    ```
    Request body:

    ```
    {
	    "user_name": "doe_j",
	    "password": "jd7845"
    }
    ```

    Response will contain a JWT token string, which has to be passed along with the headers as 'token' along with all GraphQL endpoints.

    Response:
    ```
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX25hbWUiOiJhcnVuMTIzIiwiZXhwIjoxNTM3NTg5MjU2fQ.Yy1O-BT4GTSohDU7swrPVvos761jpjwoZjjJh5qEbvI"
    ```
* GraphQL Endpoints (Method-GET)

    Set Header for all requests

    ```
    token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX25hbWUiOiJhcnVuMTIzIiwiZXhwIjoxNTM3NTg5MjU2fQ.Yy1O-BT4GTSohDU7swrPVvos761jpjwoZjjJh5qEbvI"
    ```

    * Add a post

        Request:

        ```
        http://localhost:8080/graphql?query=mutation{addPost(text:"My+new+text"){post_id,text,user_name,time_stamp}}'
        ```

        Respones:
        
        Returns the new post.

        ```
        {
            "post_id": 1,
            "text": "My new text",
            "user_name": "doe_j",
            "time_stamp": "2018-09-21"
        }
        ```

    * Follow user

        Request:

        ```
        http://localhost:8080/graphql?query=mutation{followUser(user_name:"username_to_be_folloed")}
        ```
        Respones:
        
        Returns a boolean true on successful insert.

    * Get inserted posts

        Request:

        ```
        http://localhost:8080/graphql?query={getOwnPosts{post_id,text,user_name,time_stamp}}'
        ```
        Response:

        An array of posts the user has posted.

        ```
        [
            {
                "post_id": 1,
                "text": "My new text",
                "user_name": "doe_j",
                "time_stamp": "2018-09-21"
            },
            {
                "post_id": 2,
                "text": "My other new text",
                "user_name": "doe_j",
                "time_stamp": "2018-09-22"
            }
        ]

        ```

    * Get users

        Request:

        ```
        http://localhost:8080/graphql?query={getUsers{name,user_name}}'
        ```
        Response:

        Returns an array of all users on the platform.

        ```
        [
            {
                "name": "Jane Doe",
                "user_name": "d_jane"
            },
            {
                "name": "William Smith",
                "user_name": "wsmith"
            }
        ]

        ```

    * Get feed
    
        Request:

        ```
        http://localhost:8080/graphql?query={getPostFeed{post_id,text,user_name,time_stamp}}'
        ```    

        Response:

        An array of posts by users, sorted by date, that the person follows.

        ```
        [
            {
                "post_id": 8081,
                "text": "Hello",
                "user_name": "d_jane",
                "time_stamp": "2018-09-21"
            },
            {
                "post_id": 5485,
                "text": "wsmith says Hi!",
                "user_name": "wsmith",
                "time_stamp": "2018-09-22"
            }
        ]

        ```