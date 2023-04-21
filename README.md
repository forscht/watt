# watt
Watt is an open-source SMTP wrapper written in Go that provides a simple web interface for creating and managing temporary email addresses.
It is designed to be a self-hosted solution that enables users to create and receive temporary email addresses without the need for a third-party service.

With Watt, users can easily generate a temporary email address that can be used to sign up for online services, verify accounts, or receive email notifications without revealing their real email address. The web interface provides a simple and intuitive way to manage multiple email addresses, delete old messages, and view message history.

Some of the key features of Watt include:

- Simple web interface: No bs lightning-fast web interface. Watt store recently used addresses in browser storage for forever so you don't need to remember them.

- Self-hosted: Watt is designed to be a self-hosted solution, which means that users have full control over their data and can run the service on their own servers.

- Fast and lightweight: Watt is written in Go, which makes it fast and lightweight. It is designed to handle a large number of email requests and can be easily scaled to meet the needs of high traffic websites.

- Customizable: Watt is highly customizable and can be easily configured to meet the specific needs of different users. Users can customize the SMTP server settings, TTL cache duration, and domain name.

Overall, Watt is a simple and powerful SMTP wrapper that provides a secure and convenient way to manage temporary email addresses. It is a great solution for users who want to protect their privacy online and avoid spam emails.

### Installation
1. Install Go: To use the Watt project, you first need to install the Go programming language. You can download the Go installer from the official website of Go.
2. Clone the Watt Repository: Once Go is installed, clone the Watt repository from Github using the following command in your terminal:
    ```shell
    git clone git clone https://github.com/forscht/watt.git
    ```
3. Install Dependencies: Watt has dependencies that need to be installed before it can be used. You can install them by running the following command in the cloned Watt directory:
    ```shell
    go mod tidy
    ```
4. Build watt binary with below command:
   ```shell
   go build
   ```
5. Run watt with following command:
   ```shell
   sudo ./watt --domain yourdomain.com
   ```
   

### Usage
To start Watt, simply run the built executable. Watt accepts several command-line options, which are listed below:
```shell
usage: watt --domain=DOMAIN [<flags>]

Watt: smtp wrapper for temp mail with web based interface

Flags:
  --[no-]help                Show context-sensitive help (also try --help-long and --help-man).
  --[no-]version             Show application version.
  --port=3000                Port number to start the HTTP server on ($PORT)
  --saddr=":25"              Address to start the SMTP server on ($SMTP_ADDR)
  --readtimeout=30s          Set the read timeout duration for the SMTP server ($READ_TIMEOUT)
  --writetimeout=30s         Set the write timeout duration for the SMTP server ($WRITE_TIMEOUT)
  --maxmessagebytes=1048576  Set the maximum email size in bytes for the SMTP server ($MAX_MESSAGE_BYTES)
  --domain=DOMAIN            Domain name for SMTP server. Example: 'spamok.org' ($DOMAIN)
  --name="Watt"              This will be replaced for 'Watt' in webpage ($NAME)
  --ttl=30m                  Set the time-to-live duration for the mail cache ($TTL)


```
Once Watt is running, you can generate temporary email addresses by visiting `http://localhost:3000` in your web browser.

Watt also accepts `.env` file. Below are the descriptions of each variable:

- `PORT`: The port number on which the HTTP server will listen. The default value is 3000.
- `SMTP_ADDR`: The address to start the SMTP server on. The default value is :25.
- `READ_TIMEOUT`: The read timeout duration for the SMTP server. If not set, the default value of 30 seconds will be used.
- `WRITE_TIMEOUT`: The write timeout duration for the SMTP server. If not set, the default value of 30 seconds will be used.
- `MAX_MESSAGE_BYTES`: The maximum email size in bytes for the SMTP server. The default value is 1048576 bytes (1 MB).
- `DOMAIN`: The domain name for the SMTP server. This is a required field and must be set.
- `NAME`: The name that will replace "Watt" in the web interface. The default value is "Watt".
- `TTL`: The time-to-live duration for the mail cache. The default value is 30 minutes.


### Domain records setup guide
1. Log in to your domain registrar's control panel or DNS management interface.
2. Navigate to the DNS management section for your domain.
3. Add an A record pointing to the IP address of the server where your Watt instance is running. For example, if your Watt instance is running on a server with IP address 123.45.67.89, you would create an A record with the name "watt" and the value "123.45.67.89".
4. Add an MX record pointing to the same server. The name should be "@", which is shorthand for the root domain, and the value should be the same as the A record you created in step 3, but with a priority value of 10. For example, if your A record is "watt.example.com" with value "123.45.67.89", your MX record would be "@ IN MX 10 watt.example.com.".
5. Save the changes to your DNS settings.
6. Wait for the DNS changes to propagate, which can take up to 24 hours.
7. Test your Watt instance by sending an email to an address at your domain and checking that it is received by your Watt instance.

### License
Watt is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.
