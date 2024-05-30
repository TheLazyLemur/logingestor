### Log Ingestor:

- Develop a mechanism to ingest logs in the provided format. [x]
- Ensure scalability to handle high volumes of logs efficiently. [x]
- Mitigate potential bottlenecks such as I/O operations, database write speeds, etc. [x]
- Make sure that the logs are ingested via an HTTP server, which runs on port `3000` by default. [x]

### Query Interface:

- Offer a user interface (Web UI or CLI) for full-text search across logs.
- Include filters based on:
    - level [x]
    - message [x]
    - resourceId [ ]
    - timestamp [ ]
    - traceId [x]
    - spanId [ ]
    - commit [ ]
    - metadata.parentResourceId [ ]
- Aim for efficient and quick search results.

## Advanced Features (Bonus):

These features aren’t compulsory to implement, however, adding them might increase the chances of your submission being accepted.

- Implement search within specific date ranges. [ ]
- Utilize regular expressions for search. [ ]
- Allow combining multiple filters. [ ]
- Provide real-time log ingestion and searching capabilities. [ ]
- Implement role-based access to the query interface. [ ]
