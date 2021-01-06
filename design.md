# **Linux Job Worker**

A worker service for scheduling and running linux processes. The server provides a REST API with the ability to start, stop and query job status and logs using a simple client cli.

### **Authentication**
- mTLS
- use openssl to create self signed certificates and keys.
- store cert and key files in repo folder for simplicity.  
    - obviously, would not do this in prod


### **Client**
- minimal cli for user friendliness. \
**Example:**
    - `start <linux command>`
    - `stop <job id>`
    - `list`
    - `status <job id>`
    - `logs`
- parse commands and send requests to server via json msg in req body.
- handle server response messages and communicate output to user in presentable format

### **Server**
- handle requests according to API spec below
- pass requisite info from client requests to worker for execution. 
- return response msg to client
- handle multiple clients

### **Jobs**
- handles actual execution of client requested command
    - start, stop, list, status, log
- store job properties, for example:
    - id, cmd, status, output
- manage all job related data and activities. some examples:
    - list of running jobs
    - execute linux cmds and handle results from stdout/stderr
    - provide info to server so it can send response to client
- manage resources (Concerns to be aware of. won't implement per comments)
  - track # of jobs running concurrently. dont let exceed limit.
  - cancel job if taking too long
  - limit length of command 

### **Logs**
- provide basic info about jobs. i.e. props from job struct. just prove capability. not intended to be 100% complete.
- stored in memory. will not write to file or store in DB for prototyping.
    - For production can worry about how many to keep and when to discard/no longer needed.
- may do something like only keep the 50 most recent records to limit memory consumption for now. let me know thoughts.

<br><br>

## **API**
---
### **[GET] List current Jobs**
- **GET** `/jobs` \
    Retrieve a list of current jobs (running & queued)

### **[POST] Add job to queue**
- **POST** `/jobs` \
    Add a new job to the queue

### **[GET] Get job by id**
- **GET** `/jobs/{id}`\
    Retrieve status of job matching `{id}`

### **[DELETE] Cancel job by id**
- **DELETE** `/jobs/{id}` \
    Stop job matching `{id}`

### **[DELETE] Cancel all jobs** (is this needed?)
- **DELETE** `/jobs` \
    Stop all jobs

### **[GET] Get logs**
- **GET** `/jobs/{id}/log` \
    Retrieve log for job matching `{id}` 

<br><br>

## **Limitations / Tradeoffs**
---
1. Logs stored in server-side memory. No need for persistent storage in prototype.
   - Production version of log format would be more detailed and persisted to file or db...make decisions about how many to keep/when to discard old entries.
2. Auth using self signed certs stored in repo. In production these would be stored securely. Maybe pass location in via environment variables so application can access.
3. CLI: Production version could use config file or flags for setting additional options or env's for cert/key file locations, port number, etc...
4. Just enough error handling. handle the obvious. point out where add'l attention may be required.



<br><br>
## **Questions**
---

1. Does it need functionality to cancel all jobs as well as single job by id? Or is it enough to only offer cancellation by id?
2. Re: status.  If status is complete I assume we want to present output of job as well?
3. Re: job stoppage. Probably thinking too far ahead here since restart won't be implemented but... Do we want this job to stick around in system in case you wanted to say restart by id or be fully deleted, requiring submitting an entirely new job if we wanted to restart?

I'm sure I will have more... :)