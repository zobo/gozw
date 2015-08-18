// THIS FILE IS AUTO-GENERATED BY CCGEN
// DO NOT MODIFY

package {{.CommandClass.GetPackageName}}

{{$version := .CommandClass.Version}}

func init() {
  gob.Register({{.Command.GetStructName .CommandClass}}{})
}

// {{.Help}}
type {{.Command.GetStructName .CommandClass}} struct {
  {{range $_, $param := .Command.Params}}
    {{template "command-struct-fields" $param}}
  {{end}}
}

func (cmd {{.Command.GetStructName .CommandClass}}) CommandClassID() byte {
  return {{.CommandClass.Key}}
}

func (cmd {{.Command.GetStructName .CommandClass}}) CommandID() byte {
  return byte(Command{{.Command.GetStructName .CommandClass}})
}

func (cmd *{{.Command.GetStructName .CommandClass}}) UnmarshalBinary(data []byte) error {
  // According to the docs, we must copy data if we wish to retain it after returning
  {{if .Command.Params}}
  payload := make([]byte, len(data))
  copy(payload, data)

  if len(payload) < 2 {
    return errors.New("Payload length underflow")
  }

  {{template "unmarshal-command-params" .Command.Params}}

  {{end}}
  return nil
}

func (cmd *{{.Command.GetStructName .CommandClass}}) MarshalBinary() (payload []byte, err error) {
  payload = make([]byte, 2)
  payload[0] = cmd.CommandClassID()
  payload[1] = cmd.CommandID()
  {{template "marshal-command-params" .Command.Params}}
  return
}
