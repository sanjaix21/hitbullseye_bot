# Hitbullseye Bot 🎯

An automation bot for Hitbullseye online tests that uses Google's Gemini AI to automatically answer multiple-choice questions.

## ⚠️ Disclaimer

This bot is for educational purposes only. Use it responsibly and in accordance with your institution's policies. The authors are not responsible for any misuse or violations of terms of service.

## ✨ Features

- 🤖 **AI-Powered Answers**: Uses Google Gemini AI to intelligently answer questions
- 🔄 **Fully Automated Tests**: Very little manual help is needed
- 🎯 **Anti-Detection**: Implements human-like clicking patterns to bypass anti-automation measures
- 📊 **Multiple Test Support**: Can handle multiple tests in sequence
- 🛡️ **Error Handling**: Robust error handling with fallback mechanisms

## 🛠️ Prerequisites

- **Go 1.19+** installed on your system
- **Google Chrome** browser
- **Gemini API Key** from Google

## 📦 Installation

<details>
<summary><strong>🪟 Windows</strong></summary>

1. **Install Go**:
   - Download from [https://golang.org/dl/](https://golang.org/dl/)
   - Run the installer and follow the setup wizard

2. **Install Chrome** (if not already installed):
   - Download from [https://www.google.com/chrome/](https://www.google.com/chrome/)

3. **Clone the repository**:
   ```cmd
   git clone https://github.com/sanjaix21/hitbullseye_bot.git
   cd hitbullseye_bot
   ```

</details>

<details>
<summary><strong>🍎 macOS</strong></summary>

1. **Install Go using Homebrew**:
   ```bash
   brew install go
   ```
   Or download from [https://golang.org/dl/](https://golang.org/dl/)

2. **Install Chrome**:
   ```bash
   brew install --cask google-chrome
   ```

3. **Clone the repository**:
   ```bash
   git clone https://github.com/sanjaix21/hitbullseye_bot.git
   cd hitbullseye_bot
   ```

</details>

<details>
<summary><strong>🐧 Linux</strong></summary>

### Ubuntu/Debian:
1. **Install Go**:
   ```bash
   sudo apt update
   sudo apt install golang-go
   ```

2. **Install Chrome**:
   ```bash
   wget -q -O - https://dl.google.com/linux/linux_signing_key.pub | sudo apt-key add -
   sudo sh -c 'echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google-chrome.list'
   sudo apt update
   sudo apt install google-chrome-stable
   ```

### Arch Linux:
1. **Install Go**:
   ```bash
   sudo pacman -S go
   ```

2. **Install Chrome**:
   ```bash
   # Using AUR helper (yay)
   yay -S google-chrome
   
   # Or manually from AUR
   git clone https://aur.archlinux.org/google-chrome.git
   cd google-chrome
   makepkg -si
   ```

### Clone the repository:
```bash
git clone https://github.com/sanjaix21/hitbullseye_bot.git
cd hitbullseye_bot
```

</details>

## ⚙️ Configuration

### 1. Set up Environment Variables

1. **Copy the example environment file**:
   ```bash
   cp env.example .env
   ```

2. **Edit the `.env` file** with your credentials:
   ```env
   # Hitbullseye Credentials
   HITBULLSEYE_ID=your_hitbullseye_id
   HITBULLSEYE_PASSWORD=your_password

   # Google Gemini API Key
   GEMINI_API_KEY=your_gemini_api_key_here
   ```

### 2. Get Google Gemini API Key

1. Visit [Google AI Studio](https://ai.google.dev/gemini-api/docs)
2. Sign in with your Google account
3. Create a new API key
4. Copy the API key and paste it in your `.env` file

**Detailed Guide**: [https://ai.google.dev/gemini-api/docs](https://ai.google.dev/gemini-api/docs)

### 3. Install Dependencies

```bash
go mod tidy
```

## 🚀 Usage

### Run the Bot

Simply run the following command:

```bash
go run main.go
```

The bot will automatically:
- Connect to Chrome using Rod's default connection
- Navigate to Hitbullseye
- Find and start available tests
- Use AI to answer questions
- Submit completed tests

## 🔧 Configuration Options

### Credentials Setup

Update your credentials in the `main.go` file or use environment variables:

```go
hitbullseyeId := os.Getenv("HITBULLSEYE_ID")
hitbullseyePassword := os.Getenv("HITBULLSEYE_PASSWORD")
```

### Gemini API Configuration

The bot automatically uses the `GEMINI_API_KEY` from your environment variables. If the API key is not set, it will fall back to random answers.

## 🏗️ Project Structure

```
hitbullseye_bot/
├── main.go                 # Entry point
├── internal/
│   ├── handler.go          # Main automation logic
│   └── getanswer.go        # Gemini AI integration
├── go.mod                  # Go module file
├── go.sum                  # Go dependencies
├── .env                    # Environment variables
├── env.example             # Example environment file
└── README.md               # This file
```

## 🤖 How It Works

1. **Login**: Automatically logs into Hitbullseye using provided credentials
2. **Test Navigation**: Finds and navigates to available tests
3. **Question Collection**: Scrapes all questions and options from the test
4. **AI Processing**: Sends questions to Gemini AI for intelligent answers
5. **Answer Application**: Automatically selects the AI-recommended answers
6. **Submission**: Submits the completed test

## 🛡️ Anti-Detection Features

- **Human-like Clicking**: Simulates random mouse movements and clicks
- **Timing Variations**: Uses realistic delays between actions
- **Error Handling**: Graceful handling of page load issues
- **Fallback Mechanisms**: Multiple strategies for element interaction

## 🐛 Troubleshooting

### Common Issues

1. **Chrome Connection Failed**:
   - Make sure Chrome is installed and accessible
   - Rod will automatically launch and connect to Chrome

2. **Questions Not Loading**:
   - The bot includes automatic page activation clicks
   - Manual clicking may still be required for some tests

3. **API Rate Limits**:
   - Gemini API has rate limits - the bot includes proper error handling

4. **Element Not Found Errors**:
   - Website structure may have changed
   - Check and update CSS selectors in the code

### Debug Mode

To enable verbose logging, add debug prints in the code or use Chrome DevTools to inspect elements.

## 📝 Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `HITBULLSEYE_ID` | Your Hitbullseye user ID | Yes |
| `HITBULLSEYE_PASSWORD` | Your Hitbullseye password | Yes |
| `GEMINI_API_KEY` | Google Gemini API key | Optional* |

*If not provided, the bot will use random answers as fallback.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Commit changes: `git commit -am 'Add feature'`
4. Push to branch: `git push origin feature-name`
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## ⚖️ Legal Notice

This software is provided for educational purposes only. Users are responsible for complying with their institution's academic integrity policies and the terms of service of the Hitbullseye platform. The developers do not encourage or condone academic dishonesty.

---

**Happy Testing! 🎯**
