# Chawk 

**This project is a work in progress.**<br>
A python wrapper of the blackboard api to simplify admin tasks such as managing courses.  

# Setup 
You will need to create an application to get the needed keys. You can request access to the Blackboard REST APIs through the [Developer Portal](https://developer.blackboard.com/).  


# Example of using the library 
```python

from chawk.blackboard_client import BlackboardClient


base_url = "https://blackboard.xxxx.edu"
client_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
client_secret = "xxxxxxxxxxxxxxxxxxxxxxxxx"

# Create client
client = BlackboardClient(client_id, client_secret, base_url, "example.log")

# Example usage
print(client.user.does_user_exist("00000_0"))
print(client.courses.does_course_exist("00000_1"))

client.courses.enroll_user(username="dg123", course_id="math101213", role="Student")
```

# Current Functions  

## User
 ```python

```

## Course

```python

```

## Discussion
```python

```

## Gradebook
```markdown

- Update column value for user
- Create column

```