package discord

import (
	"errors"
	"hansel/aws"
	"hansel/config"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) receive(s *discordgo.Session, event *discordgo.MessageCreate) {
	messages, err := config.GetConfig()
	if err != nil {
		log.Fatalln(err)
	}

	if event.Content == messages.StartTriggerMessage {
		// 起動時
		log.Println("開始 : インスタンス起動...")
		_, err := s.ChannelMessageSend(event.ChannelID, "インスタンスの起動コマンドを検知")
		if err != nil {
			log.Println("チャンネルメッセージの送信に失敗 :", err)
			return
		}

		err = b.awsClient.StartInstance()
		if err != nil {
			log.Println(err)

			errDiscordMsg := ""
			if errors.Is(err, aws.ErrFailedStartInstance) {
				errDiscordMsg = "インスタンスの起動に失敗"
			} else if errors.Is(err, aws.ErrInvalidResponseStartInstance) {
				errDiscordMsg = "インスタンスの起動に失敗"
			} else if errors.Is(err, aws.ErrInstanceAlreadyStarted) {
				errDiscordMsg = "インスタンスは起動済み"
			} else if errors.Is(err, aws.ErrStartingInstance) {
				errDiscordMsg = "インスタンスは既に起動準備中"
			} else if errors.Is(err, aws.ErrFailedWaitStartInstance) {
				errDiscordMsg = "インスタンスの起動状態不明　再度のコマンド入力を要求"
			}

			_, err := s.ChannelMessageSend(event.ChannelID, errDiscordMsg)
			if err != nil {
				log.Println("チャンネルメッセージの送信に失敗 :", err)
				return
			}
			return
		}

		log.Println("正常終了 : インスタンス起動")
		_, err = s.ChannelMessageSend(event.ChannelID, "インスタンスの起動に成功")
		if err != nil {
			log.Println("チャンネルメッセージの送信に失敗 :", err)
			return
		}

		// IPアドレス通知
		log.Println("IPアドレス取得待機中...")
		time.Sleep(time.Second)

		ipaddress, err := b.awsClient.GetIPAddress()
		if err != nil {
			log.Println(err)

			errDiscordMsg := "IPアドレスの取得に失敗"
			_, err := s.ChannelMessageSend(event.ChannelID, errDiscordMsg)
			if err != nil {
				log.Println("チャンネルメッセージの送信に失敗 :", err)
				return
			}
			return
		}

		_, err = s.ChannelMessageSend(event.ChannelID, "今回のIPアドレス : "+ipaddress)
		if err != nil {
			log.Println("チャンネルメッセージの送信に失敗 :", err)
			return
		}

	} else if event.Content == messages.HibernateTriggerMessage {
		// 停止時
		log.Println("開始 : インスタンス停止...")
		_, err := s.ChannelMessageSend(event.ChannelID, "インスタンスの停止コマンドを検知")
		if err != nil {
			log.Println("チャンネルメッセージの送信に失敗 :", err)
			return
		}

		err = b.awsClient.StopInstance()
		if err != nil {
			log.Println(err)

			errDiscordMsg := ""
			if errors.Is(err, aws.ErrFailedStopInstance) {
				errDiscordMsg = "インスタンスの停止に失敗"
			} else if errors.Is(err, aws.ErrInvalidResponseStopInstance) {
				errDiscordMsg = "インスタンスの停止に失敗"
			} else if errors.Is(err, aws.ErrInstanceAlreadyStopped) {
				errDiscordMsg = "インスタンスは停止済み"
			} else if errors.Is(err, aws.ErrStoppingInstance) {
				errDiscordMsg = "インスタンスは既に停止準備中"
			} else if errors.Is(err, aws.ErrFailedWaitStopInstance) {
				errDiscordMsg = "インスタンスの停止状態不明　再度のコマンド入力を要求"
			}

			_, err := s.ChannelMessageSend(event.ChannelID, errDiscordMsg)
			if err != nil {
				log.Println("チャンネルメッセージの送信に失敗 :", err)
				return
			}
			return
		}

		log.Println("正常終了 : インスタンス停止")
		_, err = s.ChannelMessageSend(event.ChannelID, "インスタンスの停止に成功")
		if err != nil {
			log.Println("チャンネルメッセージの送信に失敗 :", err)
			return
		}

	} else if event.Content == messages.GetStatusTriggerMessage {
		// 起動状態の確認(IPアドレスの取得)
		log.Println("開始 : インスタンスステータス確認")
		_, err := s.ChannelMessageSend(event.ChannelID, "インスタンスの確認コマンドを検知")
		if err != nil {
			log.Println("チャンネルメッセージの送信に失敗 :", err)
			return
		}

		ipaddress, err := b.awsClient.GetIPAddress()
		if err != nil {
			log.Println(err)

			errDiscordMsg := "インスタンスの確認に失敗"
			_, err := s.ChannelMessageSend(event.ChannelID, errDiscordMsg)
			if err != nil {
				log.Println("チャンネルメッセージの送信に失敗 :", err)
				return
			}
			return
		}

		if ipaddress != "" {
			_, err := s.ChannelMessageSend(event.ChannelID, "インスタンスは起動済み :"+ipaddress)
			if err != nil {
				log.Println("チャンネルメッセージの送信に失敗 :", err)
				return
			}
		} else {
			_, err := s.ChannelMessageSend(event.ChannelID, "インスタンスは未起動")
			if err != nil {
				log.Println("チャンネルメッセージの送信に失敗 :", err)
				return
			}
		}
	}
}
