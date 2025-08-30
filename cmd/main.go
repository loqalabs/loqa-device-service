/*
 * This file is part of Loqa (https://github.com/loqalabs/loqa).
 * Copyright (C) 2025 Loqa Labs
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program. If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/loqalabs/loqa-device-service/internal/messaging"
)

// DeviceService simulates smart home device control
type DeviceService struct {
	natsService *messaging.NATSService
	devices     map[string]*Device
}

// Device represents a smart home device
type Device struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Location string `json:"location"`
	State    string `json:"state"`
	Online   bool   `json:"online"`
}

func main() {
	log.Println("🏠 Starting Loqa Device Service")

	// Create device service
	deviceService, err := NewDeviceService()
	if err != nil {
		log.Fatalf("❌ Failed to create device service: %v", err)
	}

	// Connect to NATS
	if err := deviceService.natsService.Connect(); err != nil {
		log.Fatalf("❌ Failed to connect to NATS: %v", err)
	}

	// Start device service
	if err := deviceService.Start(); err != nil {
		log.Fatalf("❌ Failed to start device service: %v", err)
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	log.Println("🏠 Device service is running. Press Ctrl+C to stop.")
	<-sigChan

	// Shutdown
	log.Println("🛑 Shutting down device service...")
	deviceService.Shutdown()
	log.Println("👋 Device service stopped")
}

// NewDeviceService creates a new device service
func NewDeviceService() (*DeviceService, error) {
	natsService, err := messaging.NewNATSService()
	if err != nil {
		return nil, err
	}

	// Initialize some mock devices
	devices := map[string]*Device{
		"living-room-lights": {
			ID:       "living-room-lights",
			Type:     "lights",
			Name:     "Living Room Lights",
			Location: "living room",
			State:    "off",
			Online:   true,
		},
		"bedroom-lights": {
			ID:       "bedroom-lights",
			Type:     "lights", 
			Name:     "Bedroom Lights",
			Location: "bedroom",
			State:    "off",
			Online:   true,
		},
		"kitchen-lights": {
			ID:       "kitchen-lights",
			Type:     "lights",
			Name:     "Kitchen Lights", 
			Location: "kitchen",
			State:    "off",
			Online:   true,
		},
		"living-room-audio": {
			ID:       "living-room-audio",
			Type:     "audio",
			Name:     "Living Room Audio",
			Location: "living room",
			State:    "off",
			Online:   true,
		},
	}

	return &DeviceService{
		natsService: natsService,
		devices:     devices,
	}, nil
}

// Start starts the device service and subscribes to commands
func (ds *DeviceService) Start() error {
	log.Println("🔌 Starting device command subscriptions")

	// Subscribe to light commands
	_, err := ds.natsService.SubscribeToDeviceCommands("lights", ds.handleLightCommand)
	if err != nil {
		return err
	}

	// Subscribe to audio commands  
	_, err = ds.natsService.SubscribeToDeviceCommands("audio", ds.handleAudioCommand)
	if err != nil {
		return err
	}

	// Subscribe to generic device commands
	_, err = ds.natsService.SubscribeToDeviceCommands("*", ds.handleGenericCommand)
	if err != nil {
		return err
	}

	log.Printf("✅ Device service started with %d devices", len(ds.devices))
	ds.logDeviceStatus()

	return nil
}

// handleLightCommand handles light device commands
func (ds *DeviceService) handleLightCommand(event *messaging.DeviceCommandEvent) {
	log.Printf("💡 Processing light command - Action: %s, Location: %s", 
		event.Action, event.Location)

	// Find matching device
	device := ds.findDevice("lights", event.Location, event.DeviceID)
	if device == nil {
		ds.sendErrorResponse(event, "Light not found")
		return
	}

	// Execute command
	success, message := ds.executeCommand(device, event.Action)
	
	// Send response
	response := &messaging.DeviceResponseEvent{
		RequestID:  event.RequestID,
		DeviceType: "lights",
		DeviceID:   device.ID,
		Success:    success,
		Message:    message,
		Timestamp:  time.Now().UnixNano(),
	}

	if err := ds.natsService.PublishDeviceResponse(response); err != nil {
		log.Printf("❌ Failed to publish device response: %v", err)
	}
}

// handleAudioCommand handles audio device commands
func (ds *DeviceService) handleAudioCommand(event *messaging.DeviceCommandEvent) {
	log.Printf("🎵 Processing audio command - Action: %s, Location: %s", 
		event.Action, event.Location)

	// Find matching device
	device := ds.findDevice("audio", event.Location, event.DeviceID)
	if device == nil {
		ds.sendErrorResponse(event, "Audio device not found")
		return
	}

	// Execute command
	success, message := ds.executeCommand(device, event.Action)
	
	// Send response
	response := &messaging.DeviceResponseEvent{
		RequestID:  event.RequestID,
		DeviceType: "audio",
		DeviceID:   device.ID,
		Success:    success,
		Message:    message,
		Timestamp:  time.Now().UnixNano(),
	}

	if err := ds.natsService.PublishDeviceResponse(response); err != nil {
		log.Printf("❌ Failed to publish device response: %v", err)
	}
}

// handleGenericCommand handles generic device commands
func (ds *DeviceService) handleGenericCommand(event *messaging.DeviceCommandEvent) {
	log.Printf("🔧 Processing generic command - Device: %s, Action: %s", 
		event.DeviceType, event.Action)

	// This is a fallback handler for any device type
	// In a real system, you might have plugins or drivers for different device types
}

// findDevice finds a device by type, location, or ID
func (ds *DeviceService) findDevice(deviceType, location, deviceID string) *Device {
	// If device ID is specified, find by ID
	if deviceID != "" {
		if device, exists := ds.devices[deviceID]; exists && device.Type == deviceType {
			return device
		}
	}

	// Otherwise, find by type and location
	for _, device := range ds.devices {
		if device.Type == deviceType {
			if location == "" || device.Location == location {
				return device
			}
		}
	}

	return nil
}

// executeCommand executes a command on a device
func (ds *DeviceService) executeCommand(device *Device, action string) (bool, string) {
	if !device.Online {
		return false, fmt.Sprintf("%s is offline", device.Name)
	}

	switch action {
	case "on":
		if device.State == "on" {
			return true, fmt.Sprintf("%s is already on", device.Name)
		}
		device.State = "on"
		return true, fmt.Sprintf("%s turned on", device.Name)

	case "off":
		if device.State == "off" {
			return true, fmt.Sprintf("%s is already off", device.Name)
		}
		device.State = "off"
		return true, fmt.Sprintf("%s turned off", device.Name)

	case "play":
		if device.Type == "audio" {
			device.State = "playing"
			return true, fmt.Sprintf("%s started playing", device.Name)
		}
		return false, fmt.Sprintf("Cannot play on %s", device.Name)

	case "stop":
		if device.Type == "audio" {
			device.State = "stopped"
			return true, fmt.Sprintf("%s stopped", device.Name)
		}
		return false, fmt.Sprintf("Cannot stop %s", device.Name)

	case "pause":
		if device.Type == "audio" {
			device.State = "paused"
			return true, fmt.Sprintf("%s paused", device.Name)
		}
		return false, fmt.Sprintf("Cannot pause %s", device.Name)

	default:
		return false, fmt.Sprintf("Unknown action: %s", action)
	}
}

// sendErrorResponse sends an error response
func (ds *DeviceService) sendErrorResponse(event *messaging.DeviceCommandEvent, message string) {
	response := &messaging.DeviceResponseEvent{
		RequestID:  event.RequestID,
		DeviceType: event.DeviceType,
		DeviceID:   event.DeviceID,
		Success:    false,
		Message:    message,
		Timestamp:  time.Now().UnixNano(),
	}

	if err := ds.natsService.PublishDeviceResponse(response); err != nil {
		log.Printf("❌ Failed to publish error response: %v", err)
	}
}

// logDeviceStatus logs the status of all devices
func (ds *DeviceService) logDeviceStatus() {
	log.Println("📊 Device Status:")
	for _, device := range ds.devices {
		status := "🔴 offline"
		if device.Online {
			if device.State == "on" || device.State == "playing" {
				status = "🟢 " + device.State
			} else {
				status = "🟡 " + device.State
			}
		}
		log.Printf("  %s (%s in %s): %s", device.Name, device.Type, device.Location, status)
	}
}

// Shutdown shuts down the device service
func (ds *DeviceService) Shutdown() {
	if ds.natsService != nil {
		ds.natsService.Close()
	}
}