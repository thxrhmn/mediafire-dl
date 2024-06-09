# MediaFire-dl

MediaFire-dl is a command-line tool that allows you to download files from MediaFire URLs with a progress bar.

## Installation

1. Ensure you have Go installed on your machine. You can download it [here](https://go.dev/dl/).

2. Clone the repository:
    ```bash
    git clone https://github.com/thxrhmn/mediafire-dl.git
    ```

3. Change into the project directory:
    ```bash
    cd mediafire-dl
    ```

4. Install the project dependencies:
    ```bash
    go mod tidy
    ```

5. Build the project:
    ```bash
    go build -o mediafire-dl .

    ```

6. Optionally, install the binary globally on your system:

    ```bash
    go install
    ```

7. Run the program 
    ```bash
    ./mediafire-downloader [flags] [URLs]

    ```

## Usage
To download from multiple URLs provided as arguments:
```bash
./mediafire-downloader "url1" "url2"
```

To download from URLs listed in a file:
```bash
./mediafire-downloader -f urls.txt
```