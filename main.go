package main

import (
  "fmt"
  "log"
  "os"
  "os/exec"
  "regexp"
  "runtime"
  "strconv"
  "strings"

  tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


const HELP_TEXT = `Here's some help to get you started.

* setvol xx%: Set volume to xx percent.
* sleep: Put the comptuer into sleep mode
* help: Print the help text
`


func setVolume(percentage float64) bool {
  if runtime.GOOS == "windows" {
    vol := int((percentage / 100.0) * 65535.0)
    cmd := exec.Command("nircmd", "setsysvolume", strconv.Itoa(vol))
    err := cmd.Run()
    return err == nil
  } else if runtime.GOOS == "linux" {
    cmd := exec.Command("amixer", "sset", "Master", fmt.Sprintf("%f%%", percentage))
    err := cmd.Run()
    return err == nil
  }
  return false
}


func goToSleep() bool {
  if runtime.GOOS == "windows" {
    cmd := exec.Command("nircmd", "standby")
    err := cmd.Run()
    return err == nil
  } else if runtime.GOOS == "linux" {
    return false // TODO
  }
  return false
}


func trySetvolCommand(messageText string) string {
    r := regexp.MustCompile(`setvol ([0-9]+)%`)
    submatch := r.FindStringSubmatch(messageText)
    if len(submatch) != 2 {
      return "Wait, what do you mean? (usage: setvol xx%)"
    }

    percentage, err := strconv.Atoi(submatch[1])
    if err != nil {
      return "Wait, what do you mean? (usage: setvol xx%)"
    }

    if setVolume(float64(percentage)) {
      return fmt.Sprintf("Okay, set volume to %d%%", percentage)
    } else {
      return "That didn't work (error running volume command)"
    }
}


func trySleepCommand() string {
  if goToSleep() {
    return "Nighty night (sleeping)"
  } else {
    return "That didn't work (error running sleep command)"
  }
}


func tryHelpCommand() string {
  return HELP_TEXT
}


func processMessage(update tgbotapi.Update) string {
  if strings.HasPrefix(update.Message.Text, "setvol") {
    return trySetvolCommand(update.Message.Text)
  } else if strings.HasPrefix(update.Message.Text, "sleep") {
    return trySleepCommand()
  } else if strings.HasPrefix(update.Message.Text, "help") {
    return tryHelpCommand()
  }

  return "Whatever you say :) (command not found)"
}


func main() {
  bot, err := tgbotapi.NewBotAPI(os.Getenv("ALICE_TOKEN"))
  if err != nil {
    log.Panic(err)
  }

  bot.Debug = true

  log.Printf("Authorized on account %s", bot.Self.UserName)

  u := tgbotapi.NewUpdate(0)
  u.Timeout = 60

  updates := bot.GetUpdatesChan(u)

  for update := range updates {
    if update.Message != nil {
      replyText := processMessage(update)
      msg := tgbotapi.NewMessage(update.Message.Chat.ID, replyText)
      msg.ReplyToMessageID = update.Message.MessageID
      bot.Send(msg)
    }
  }
}
