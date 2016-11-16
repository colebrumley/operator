package operator

// DefaultConfig is the default server configuration that is loaded then overwritten with custom config options.
var DefaultConfig = map[string]interface{}{
	"broker": map[string]interface{}{
		// "url": "amqp://guest:guest@localhost:5672/",
		"url": "redis://localhost:6379",
		"exchange": map[string]interface{}{
			"name": "operator",
			"type": "direct",
		},
		"queue":      "op_tasks",
		"bindingkey": "op_task",
	},
	"log": map[string]interface{}{
		"level":       "info",
		"destination": "stdout",
	},
	"results": map[string]interface{}{
		"url": "redis://localhost:6379",
	},
	"tasks": map[string]interface{}{
		"enabled": []interface{}{
			"ping",
			"ssh",
			"exec",
		},
	},
	"api": map[string]interface{}{
		"enabled":   true,
		"addr":      ":8080",
		"basicauth": false,
		"usetls":    false,
		"password":  "changeme",
	},
}
