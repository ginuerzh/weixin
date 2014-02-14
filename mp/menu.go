// menu
package mp

import ()

type Button struct {
	Type      string   `json:"type"`
	Name      string   `json:"name"`
	Key       string   `json:"key,omitempty"`
	Url       string   `json:"url,omitempty"`
	SubButton []Button `json:"sub_button,omitempty"`
}

func (button *Button) AddSubButton(btn Button) {
	if len(button.SubButton) == 5 {
		return
	}
	button.SubButton = append(button.SubButton, btn)
}

type ButtonList struct {
	Buttons []Button `json:"button"`
}

type Menu struct {
	buttons ButtonList `json:"menu"`
	Error
}

func NewMenu() *Menu {
	return &Menu{}
}

func (menu *Menu) Size() int {
	return len(menu.buttons.Buttons)
}

func (menu *Menu) AddClickButton(name, key string) {
	if menu.Size() == 3 {
		return
	}
	btn := Button{Type: "click", Name: name, Key: key}
	menu.buttons.Buttons = append(menu.buttons.Buttons, btn)
}

func (menu *Menu) AddViewButton(name, url string) {
	if menu.Size() == 3 {
		return
	}
	btn := Button{Type: "view", Name: name, Url: url}
	menu.buttons.Buttons = append(menu.buttons.Buttons, btn)
}

func (menu *Menu) AddClickSubButton(index int, name, key string) {
	if index < 0 || index >= menu.Size() {
		return
	}
	btn := Button{Type: "click", Name: name, Key: key}
	menu.buttons.Buttons[index].AddSubButton(btn)
}

func (menu *Menu) AddViewSubButton(index int, name, url string) {
	if index < 0 || index >= menu.Size() {
		return
	}
	btn := Button{Type: "view", Name: name, Url: url}
	menu.buttons.Buttons[index].AddSubButton(btn)
}
