package main

import (
	"fmt"
	"strings"
	"runtime"
	"time"
	"io/ioutil"
	"os"

	"./linego/LineThrift"
	"./linego/auth"
	"./linego/helper"
	"./linego/service"
	"./linego/talk"
	"./linego/config"
)

var argsRaw = os.Args
var Basename = argsRaw[1]
var ArgSname = argsRaw[2]
var AppName = argsRaw[3]

//type User struct {
//    Group string
//    Link int
//    Invite int
//    Kick int
//}
//var Protection = []User{}
var  ProQR = []string{}
var  ProInvite = []string{}
var  ProKick = []string{}


func checkEqual(list1 []string, list2 []string) bool {
	for _, v := range list1 {
		if helper.InArray(list2, v) {
			return true
		}
	}
	return false
}

func canceling(group string, korban []string) {
					runtime.GOMAXPROCS(10)
					for _, vo := range korban {
						cancel := []string{vo}
						if helper.InArray(service.Banned, vo) {
							go func() {
								talk.CancelGroupInvitation(group, cancel)
							}()
						}
					}
}

func addbl(pelaku string) {
		if !helper.IsBanned(pelaku) {
				service.Banned = append(service.Banned, pelaku)
		}
}


func cancelall(group string, korban []string) {
					runtime.GOMAXPROCS(10)
					for _, vo := range korban {
						    cancel := []string{vo}
							go func() {
								talk.CancelGroupInvitation(group, cancel)
							}()
					}
}

func checkurl(group string, pelaku string) {
			runtime.GOMAXPROCS(10)
			go func() {
				res, _ := talk.GetGroup(group)
				cek := res.PreventedJoinByTicket
				if !cek {
					go func() {
						res.PreventedJoinByTicket = true
						talk.UpdateGroup(res)
					}()
					go func() {
						talk.KickoutFromGroup(group, []string{pelaku})
					}()
				}
			}()
}


func bot(op *LineThrift.Operation) {
	var Mid string = service.MID

	if op.Type == 26 {
		msg := op.Message
		sender := msg.From_
		var sname = ArgSname
		var rname = Basename
		var txt string
		var pesan = strings.ToLower(msg.Text)

		if strings.HasPrefix(pesan , rname + " ") {
			txt = strings.Replace(pesan, rname + " ", "", 1)
		} else if strings.HasPrefix(pesan , rname) {
			txt = strings.Replace(pesan, rname, "", 1)
		} else if strings.HasPrefix(pesan , sname + " ") {
			txt = strings.Replace(pesan, sname + " ", "", 1)
		} else if strings.HasPrefix(pesan , sname){
			txt = strings.Replace(pesan, sname, "", 1)
		}

		if sender != "" && helper.IsAccess(msg.From_) {
			if strings.Contains(pesan , "kick") {
            	str := fmt.Sprintf("%v",msg.ContentMetadata["MENTION"])
            	taglist := helper.GetMidFromMentionees(str)
            	if taglist != nil {
            		for _,target := range taglist {
            			if !helper.IsBanned(target) {
            				addbl(target)
            			}
            		}
            	}
            }
			if txt == "speed" {
				start := time.Now()
				talk.SendMessage(msg.To, "Tes....", 2)
				elapsed := time.Since(start)
				stringTime := elapsed.String()
				talk.SendMessage(msg.To, stringTime, 2)
			} else if pesan == "res" {
				talk.SendMessage(msg.To, rname, 2)
			} else if pesan == "sname" {
				talk.SendMessage(msg.To, sname, 2)
			} else if txt == "stafflist"{
				nm := []string{}
				for c, a := range service.Creator {
					res,_ := talk.GetContact(a)
					name := res.DisplayName
					c += 1
					name = fmt.Sprintf("%v. %s",c , name)
					nm = append(nm, name)
				}
				stf := "List staff:\n\n"
				str := strings.Join(nm, "\n")
				talk.SendMessage(msg.To, stf+str, 2)
			} else if txt == "status" {
					anu := talk.InviteIntoGroup(msg.To, service.Creator)
					if anu != nil {
						talk.SendMessage(msg.To, "Limit", 2)
					} else {
						talk.SendMessage(msg.To, "Normal", 2)
					}
			} else if txt == "banlist"{
				nm := []string{}
				for c, a := range service.Banned {
					res,_ := talk.GetContact(a)
					name := res.DisplayName
					c += 1
					name = fmt.Sprintf("%v. %s",c , name)
					nm = append(nm, name)
				}
				stf := "Shitlist:\n\n"
				str := strings.Join(nm, "\n")
				talk.SendMessage(msg.To, stf+str, 2)
			} else if txt == "squadlist"{
				nm := []string{}
				for c, a := range service.Squad {
					res,_ := talk.GetContact(a)
					name := res.DisplayName
					c += 1
					name = fmt.Sprintf("%v. %s",c , name)
					nm = append(nm, name)
				}
				stf := "Squad list:\n\n"
				str := strings.Join(nm, "\n")
				talk.SendMessage(msg.To, stf+str, 2)
			} else if txt == "nukeall" {
				runtime.GOMAXPROCS(100)
				res, _ := talk.GetGroup(msg.To)
				memlist := res.Members
				for _, v := range memlist {
					if !helper.IsAccess(v.Mid) {
						cancel := []string{v.Mid}
						go func() {
							talk.KickoutFromGroup(msg.To, cancel)
						}()
					}
				}
			} else if txt == "cancelall" {
				runtime.GOMAXPROCS(100)
				res, _ := talk.GetGroup(msg.To)
				memlist := res.Invitee
				for _, v := range memlist {
					if !helper.IsAccess(v.Mid) {
						cancel := []string{v.Mid}
						go func() {
							talk.CancelGroupInvitation(msg.To, cancel)
						}()
					}
				}
			} else if txt == "bye" {
				talk.LeaveGroup(msg.To)
			} else if txt == "clearban" {
				jum := len(service.Banned)
				service.Banned = []string{}
				str := fmt.Sprintf("Cleared %v banlist.", jum)
				talk.SendMessage(msg.To, str, 2)
			} else if txt == "invitebot" {
				    talk.InviteIntoGroup(msg.To, service.Squad)
			} else if txt == "open" {
				res, _ := talk.GetGroup(msg.To)
				cek := res.PreventedJoinByTicket
				if cek {
					res.PreventedJoinByTicket = false
					talk.UpdateGroup(res)
					fmt.Println("done")
				} else {
					res.PreventedJoinByTicket = true
					talk.UpdateGroup(res)
					fmt.Println("done")
			    }
            } else if strings.HasPrefix(txt , "kick") {
            	str := fmt.Sprintf("%v",msg.ContentMetadata["MENTION"])
            	taglist := helper.GetMidFromMentionees(str)
            	if taglist != nil {
            		for _,target := range taglist {
            			runtime.GOMAXPROCS(10)
            			go func() {
            				talk.KickoutFromGroup(msg.To, []string{target})
            			}()
            		}
            	}
            } else if strings.HasPrefix(txt , "addstaff") {
            	str := fmt.Sprintf("%v",msg.ContentMetadata["MENTION"])
            	taglist := helper.GetMidFromMentionees(str)
            	if taglist != nil {
            		lisa := []string{}
            		for c,target := range taglist {
            			if !helper.InArray(service.Creator, target) {
            				service.Creator = append(service.Creator, target)
            			}
            			res,_ := talk.GetContact(target)
					    name := res.DisplayName
					    c += 1
						name = fmt.Sprintf("%v. %s",c , name)
						lisa = append(lisa, name)
            		}
            		stf := "Added to staff:\n\n"
					str = strings.Join(lisa, "\n")
					talk.SendMessage(msg.To, stf+str, 2)
            	}

            } else if strings.HasPrefix(txt , "addban") {
            	str := fmt.Sprintf("%v",msg.ContentMetadata["MENTION"])
            	taglist := helper.GetMidFromMentionees(str)
            	if taglist != nil {
            		lisa := []string{}
            		for c,target := range taglist {
            			if !helper.IsBanned(target) {
            				addbl(target)
            			}
            			res,_ := talk.GetContact(target)
					    name := res.DisplayName
					    c += 1
						name = fmt.Sprintf("%v. %s",c , name)
						lisa = append(lisa, name)
            		}
            		stf := "Added to Shitlist:\n\n"
					str := strings.Join(lisa, "\n")
					talk.SendMessage(msg.To, stf+str, 2)
            	}
            } else if txt == "proqr on" {
            	if !helper.InArray(ProQR, msg.To) {
            		ProQR = append(ProQR, msg.To)
            		talk.SendMessage(msg.To, "Link Protection Enabled.", 2)
            	} else {
            		talk.SendMessage(msg.To, "Link Protection Already Enabled.", 2)
            	}
			} else if txt == "proinvite on" {
            	if !helper.InArray(ProInvite, msg.To) {
            		ProInvite = append(ProInvite, msg.To)
            		talk.SendMessage(msg.To, "Invitation protection enabled.", 2)
            	} else {
            		talk.SendMessage(msg.To, "Invitation protection already Enabled.", 2)
            	}
            } else if txt == "prokick on" {
            	if !helper.InArray(ProKick, msg.To) {
            		ProKick = append(ProKick, msg.To)
            		talk.SendMessage(msg.To, "Kick protection enabled.", 2)
            	} else {
            		talk.SendMessage(msg.To, "Kick protection already Enabled.", 2)
            	}
            } else if txt == "proqr off" {
            	if helper.InArray(ProQR, msg.To) {
            		ProQR = helper.Remove(ProQR, msg.To)
            		talk.SendMessage(msg.To, "Link protection disabled.", 2)
            	} else {
            		talk.SendMessage(msg.To, "Link protection already disabled.", 2)
            	}
			} else if txt == "proInvite off" {
            	if helper.InArray(ProInvite, msg.To) {
            		ProInvite = helper.Remove(ProInvite, msg.To)
            		talk.SendMessage(msg.To, "Invitation protection disabled.", 2)
            	} else {
            		talk.SendMessage(msg.To, "Invitation protection already disabled.", 2)
            	}
			} else if txt == "prokick off" {
            	if helper.InArray(ProKick, msg.To) {
            		ProQR = helper.Remove(ProKick, msg.To)
            		talk.SendMessage(msg.To, "Kick protection disabled.", 2)
            	} else {
            		talk.SendMessage(msg.To, "Kick protection already disabled.", 2)
            	}
			} else if txt == "set" {
				checking := []string{}
				stf := "Bot Setting:\n\n"
				if helper.InArray(ProKick, msg.To) {
					na := "༗ Kick Protect  ⇉ On"
					checking = append(checking, na)
				} else {
					na := "༗ Kick Protect  ⇉ Off"
					checking = append(checking, na)
				}
				if helper.InArray(ProInvite, msg.To) {
					na := "༗ Deny Invite   ⇉ On"
					checking = append(checking, na)
				} else {
					na := "༗ Deny Invite   ⇉ Off"
					checking = append(checking, na)
				}
				if helper.InArray(ProQR, msg.To) {
					na := "༗ Link Protect  ⇉ On"
					checking = append(checking, na)
				} else {
					na := "༗ Link Protect  ⇉ Off"
					checking = append(checking, na)
				}
				str := strings.Join(checking, "\n")
				talk.SendMessage(msg.To, stf+str, 2)
			}
		}
	} else if op.Type == 19 {
		korban := op.Param3
		kicker := op.Param2
		group := op.Param1

		if korban == Mid {
            go func() {
            	addbl(kicker)
            }()
		} else if helper.InArray(service.Squad, korban) {
			runtime.GOMAXPROCS(10)
			go func() {
				talk.KickoutFromGroup(group, []string{kicker})
			}()
			go func() {
				talk.InviteIntoGroup(group, service.Squad)
			}()
			go func() {
            	addbl(kicker)
            }()
		} else if helper.IsAccess(korban) {
			runtime.GOMAXPROCS(10)
			go func() {
				talk.KickoutFromGroup(group, []string{kicker})
			}()
			go func() {
				talk.InviteIntoGroup(group, []string{korban})
			}()
			go func() {
            	addbl(kicker)
            }()
		} else if helper.InArray(ProKick, group) {
			go func(){
				talk.KickoutFromGroup(group, []string{kicker})
			}()
			go func() {
            	addbl(kicker)
            }()
		}

	} else if op.Type == 13 {
		runtime.GOMAXPROCS(30)
	    korban := strings.Split(op.Param3, "\x1e")
		inviter := op.Param2
		group := op.Param1

		if helper.InArray(korban, Mid) {
			if helper.IsAccess(inviter) {
				go func() {
					talk.AcceptGroupInvitation(group)
				}()
			}

		} else if helper.InArray(ProInvite, group) {
			if !helper.IsAccess(inviter) {
					go func() {
						cancelall(group, korban)
					}()
					go func() {
						talk.KickoutFromGroup(group, []string{inviter})
					}()
					go func() {
            			addbl(inviter)
           			}()
			}
		} else if checkEqual(korban, service.Banned) {
			if !helper.IsAccess(inviter) {
					go func() {
						canceling(group, korban)
					}()
					go func() {
						talk.KickoutFromGroup(group, []string{inviter})
					}()
					go func() {
            				addbl(inviter)
            		}()
			}
		} else if helper.IsBanned(inviter) {
					go func() {
						cancelall(group, korban)
					}()
					go func() {
						talk.KickoutFromGroup(group, []string{inviter})
					}()
		}

	} else if op.Type == 32 {
		runtime.GOMAXPROCS(20)
		korban := op.Param3
		kicker := op.Param2
		group := op.Param1

		if helper.IsAccess(korban) && !helper.IsAccess(kicker) {
			go func(){
				talk.KickoutFromGroup(group, []string{kicker})
			}()
			go func(){
				talk.InviteIntoGroup(group, service.Squad)
			}()
			go func() {
            	addbl(kicker)
            }()
		} else if helper.InArray(ProKick, group) {
			go func(){
				talk.KickoutFromGroup(group, []string{kicker})
			}()
			go func() {
            	addbl(kicker)
            }()
		}
	} else if op.Type == 17 {
		runtime.GOMAXPROCS(10)
		kicker := op.Param2
		group := op.Param1

		if helper.IsBanned(kicker) {
			go func(){
				talk.KickoutFromGroup(group, []string{kicker})
			}()
		}
	} else if op.Type == 11 {
		runtime.GOMAXPROCS(10)
		changer := op.Param2
		group := op.Param1

		if helper.InArray(ProQR, group) {
			go func(){
				checkurl(group, changer)
			}()
		} else if helper.IsBanned(changer) {
			go func(){
				checkurl(group, changer)
			}()
		} 
	}
}


func main() {

	filepath := fmt.Sprintf("/root/gogo/token/%s.txt", Basename)
    b, err := ioutil.ReadFile(filepath)
    if err != nil {
        fmt.Print(err)
    }
    token := string(b)
    config.LINE_APPLICATION = AppName

    be, err := ioutil.ReadFile("/root/gogo/squad.txt")
    if err != nil {
        fmt.Print(err)
    }
    sqd := string(be)
    squad := strings.Split(sqd, ",")
    for _,sq := range squad {
    	service.Squad = append(service.Squad, sq)
    }
    fmt.Println(service.Squad)
    
	auth.LoginWithAuthToken(token)
	//auth.LoginWithQrCode(true)
	for {
		fetch, _ := talk.FetchOperations(service.Revision, 1)
		if len(fetch) > 0 {
			//if error != nil {
			//	fmt.Println(error)
			//}
			rev := fetch[0].Revision
			service.Revision = helper.MaxRevision(service.Revision, rev)
			bot(fetch[0])
		}
	}
}