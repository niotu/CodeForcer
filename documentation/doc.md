# Documentation

<details>
 <summary><code>POST</code><code><b>/setAdmin</b></code> <code>(sets admin's confidential data)</code></summary>

##### Parameters

> | name     | type     | data type |
> |----------|----------|-----------|
> | key      | required | string    |
> | secret   | required | string    |
> | handle   | required | string    |
> | password | required |  string   |

##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `201`         | `text/plain;charset=UTF-8`        | `Configuration created successfully`                                |
> | `400`         | `application/json`                | `{"code":"400","message":"Bad Request"}`                            |
> | `405`         | `text/html;charset=utf-8`         | None                                                                |

</details>

------------------------------------------------------------------------------------------

<details>
 <summary><code>POST</code> <code><b>/getGroups</b></code> <code>(sets admin's confidential data)</code></summary>

##### Parameters

> None

</details>

------------------------------------------------------------------------------------------

<details>
 <summary><code>POST</code> <code><b>/getContests</b></code> <code>(sets admin's confidential data)</code></summary>

##### Parameters

> | name     | type     | data type |
> |----------|----------|-------|
> | groupCode| required | string|


</details>

------------------------------------------------------------------------------------------

<details>
 <summary><code>POST</code> <code><b>/proceed</b></code> <code>(starts the process of parsing points and solutions)</code></summary>

##### Parameters

> | name      | type     | data type |
> |-----------|----------|-----------|
> | groupCode | required | string    |
> | contestID | required | int       |

</details>

------------------------------------------------------------------------------------------
