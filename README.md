# xmlpull
A simple pull parser "adapter" for golang encoding/xml

```go
	parser := xmlpull.NewParserBytes(request_data)
	atoms := parser.GetAtoms()

	a_ns1 := atoms.AddAtom("http://schemas.../ns1")
	a_ns2 := atoms.AddAtom("http://schemas.../ns2")

	a_notification := atoms.AddAtom("Notification")
	a_more_events := atoms.AddAtom("MoreEvents")

	for {
		// Read tokens from the XML document in a stream.
		t, err := parser.NextToken()
		if err != nil {
			return err
		}

		if t == nil {
			break
		}

		switch el := t.(type) {
		case xmlpull.Tag:
			// debug
			if el.IsStart {
				fmt.Println(">", el.Name, el.Space, el.Local)
			}
			if el.IsEnd {
				fmt.Println("<", el.Name, el.Space, el.Local)
			}
			
			// usage
			if el.IsTag(a_ns1, a_notification) {
				if el.IsStart {
					// ... we're entering ns1:Notification
				}
				if el.IsEnd {
					// ... we're leaving ns1:Notification
				}
			}

		case xmlpull.Text:
		  // debug
			fmt.Println(el.Text)
			
			// usage
			if el.Tag.isTag(a_ns2, a_more_events) {
				// ... get text from inside ns2:MoreEvents
				more := el.AsBoolean()	
			}
		}
	}
```
