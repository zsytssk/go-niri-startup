# niri-go-startup

this is a go script for some niri wm functions. eg `scratchpad` `switch screen` and `run app or focus existed app` ..

## To install dependencies:

```bash
go build
```

To run:

```bash
niri-startup
```

you can run this script in niri config

```
spawn-sh-at-startup "niri-startup"
```

## bind for niri shortcut

```
Alt+5  hotkey-overlay-title="Run spad chatgpt" {  spawn-sh "curl -X POST http://127.0.0.1:6321/spad -d '{\"name\":\"chatgpt\"}'"; }
Mod+T hotkey-overlay-title="Open a Terminal: ghostty" {  spawn-sh "curl -X POST http://127.0.0.1:6321/runApp -d '{\"app_id\":\"com.mitchellh.ghostty\",\"cmd\":\"ghostty\"}'"; }
Mod+Shift+Ctrl+K { spawn-sh "curl -X POST http://127.0.0.1:6321/action -d '{\"name\":\"switch-screen-prev\"}'"; }
Mod+Shift+Ctrl+J { spawn-sh "curl -X POST http://127.0.0.1:6321/action -d '{\"name\":\"switch-screen-next\"}'"; }
```
