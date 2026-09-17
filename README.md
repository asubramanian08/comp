# Fall 2026

## Step 1: Introduction

You will create a project from scratch using the latest version of SvelteKit for the frontend and one of the following backend options:

- Go with a JSON API
- Go with a GraphQL API
- Python with an HTTP API

The objective is to become familiar with how frontend and backend tools work together and to practice learning efficiently from documentation. The instructions are intentionally sparse, so you should expect to research some implementation details independently.

Use the official documentation for SvelteKit and for your chosen backend language, framework, and API style whenever possible.

### Notes

- Svelte extends HTML with additional syntax and compiles into HTML, CSS, and JavaScript that run in the browser.
- SvelteKit is a framework for building web applications with Svelte. In this project, it will provide the browser-based frontend.
- Go is a statically typed, garbage-collected language that compiles into a binary and provides straightforward concurrency tools.
- Python is a dynamically typed language with web frameworks such as FastAPI and Flask.
- JSON is a text-based format commonly used to exchange structured data between a client and server.
- A JSON API exposes HTTP endpoints that accept requests and return JSON responses.
- GraphQL is an API query language that allows a client to request and modify data through a typed schema.

The [OSI model](https://en.wikipedia.org/wiki/OSI_model) may help you understand how the networking layers relate. For this project, the browser and backend communicate using HTTP over TCP/IP. In production, this would usually use HTTPS, meaning the HTTP communication is protected by TLS.

Your application will contain two main parts:

1. A SvelteKit frontend that runs in the browser.
2. A Go or Python backend that stores and updates a shared counter.

When the user clicks a button, the SvelteKit frontend will send an HTTP request to the backend. The backend will update the counter and return its new value through the API.

**Think about and research why these tools are used, how they work together, and the advantages and disadvantages of your chosen language and API design.** Understanding the system is more important than adding extra features or elaborate styling.

## Step 2: Shared Counter

Create a button in the browser. Whenever a user clicks it:

1. The frontend sends a request to the backend.
2. The backend increments a shared server-side counter by 1.
3. The backend returns the updated value.
4. The frontend displays the returned value.

The counter must be shared among all connected clients. For example, suppose one client displays `5`. If other clients click the button before that client clicks again, the next value it receives might be `9` rather than `6`.

The backend must safely handle requests from multiple clients at the same time. If you use Go, research atomics or mutexes. If you use Python, research how your chosen server handles concurrent requests and how shared state should be protected.

There are no additional requirements for style or functionality. Keep the project simple and focus on understanding how the frontend, API, and backend work together.

### Questions to Consider

- Which backend option will you use, and why?
- What endpoint or GraphQL mutation should the frontend call?
- Which HTTP method should a JSON API use for the increment operation?
- What data, if any, should the request contain?
- What should the response contain?
- What numeric type should store the counter?
- How will the backend prevent concurrent requests from interfering with one another?
- How will the SvelteKit frontend communicate with the backend?
- Will the browser requests require CORS configuration?
- What should the frontend do if a request fails?
- What limitations result from storing the counter only in server memory?
