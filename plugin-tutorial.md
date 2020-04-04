---
layout: default
title: Plugin Tutorial
permanlink: /plugins/tutorial/
weight: 3
---

## Before we start

This is a step-by-step tutorial that shows how to build a QuestScreen plugin.
It deliberately does not describe the design of the plugin API in detail; this is what the [Plugin Documentation](/plugins/documentation/) is for.

This article assumes you are somewhat familiar with Go, JavaScript, HTML and programming in general.
It tries to describe everything good enough so that a beginner should be able to understand it.

As for the build environment, you'll need the [Go compiler][1] (at least Go 1.12), the [SDL2 library][2] including its headers, on a supported operating system (Linux and macOS are both fine, Windows is unsupported).

### What we're gonna do

In this tutorial, we'll be writing a module which displays an in-universe calendar.
As lots of settings have custom in-universe calendars, this seems to be a worthy endeavor which can be a prototype for other calendars.
The calendar we will implement is the one of [Discworld][3].

The display should show the current day, month and year.
We'll be counting years according to the University Calendar, which counts *common years* having 400 days each.

The user should be able to skip through the calendar easily.
Therefore, we want to have **+1**, **+5**, **-1** and **-5** buttons for days, months and years.
An additional editor for selecting a specific date is left as an exercise for the reader.

## Creating the skeleton

We'll be starting from the [plugin template](https://github.com/QuestScreen/PluginTemplate) repository.
For a real plugin, you might want to immediately create your own repository using the *Use this template* button, but for now, we'll simply be downloading the contents as ZIP.
Unzip it into a folder `discworld` so that this folder contains the file `go.mod`.

The directory `moduleTemplate` contains the template for our module.
A plugin may contain any number of modules, so for each additional module, you'd just copy the module template.
We'll rename the folder to `calendar`.

## Implementing the module

### Calendar data

First, we need to think about which data needs to be stored by our module.
As we want to display a calendar, we'll need to store a date.
Now thankfully, unlike the Gregorian calendar, Ankh-Morpork's calendar has is very regular – no leap years etc.
So let's implement a simple type that stores a data of the Ankh-Morpork calendar.
Put this into a file `universitydate.go` within the `calendar` directory:

```Go
package calendar

import "encoding/json"

// Month holds a Discworld month.
type Month int

// Go has no enums, but syntax that allows defining constants as
// consecutive numbers.
const (
	Ick Month = iota
	Offle
	February
	March
	April
	May
	June
	Grune
	August
	Spune
	Sektober
	Ember
	December
)

func (m Month) String() string {
	return [...]string{"Ick", "Offle", "February", "March", "April", "May",
		"June", "Grune", "August", "Spune", "Sektober", "Ember", "December"}[m]
}

// UniversityDate stores days since the 1st of Ick, year 0.
// for simplicity, we assume there is a year 0;
// I can't find official information on it.
type UniversityDate int

func (d UniversityDate) year() int {
	return int(d / 400) // a common year has 400 days.
}

func (d UniversityDate) month() Month {
	// Ick has 16 days, all other months have 32 days.
	return Month((d%400 + 16) / 32)
}

func (d UniversityDate) dayOfMonth() int {
	t := d % 400
	if t < 16 {
		// 1-based, therefore +1
		return int(t + 1)
	}
	return int(t+16)/32 + 1
}

func (d UniversityDate) add(daysDelta int) UniversityDate {
	return d + UniversityDate(daysDelta)
}

// MarshalJSON returns an object containing Day, Month and Year.
func (d UniversityDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Day   int
		Month string
		Year  int
	}{
		Day: d.dayOfMonth(), Month: d.month().String(), Year: d.year(),
	})
}
```

As you can see, we defined `MarshalJSON` to customize the way this type is serialized to JSON.
We used an anonymous `struct` type that contains day, month and year instead of the total days we store internally.

QuestScreen serializes data in two ways:
Everything being sent to and received from the Web Client is serialized as JSON.
Everything being stored to and loaded from the file system is serialized as YAML.
This makes it easy for you to define the way your data is serialized/deserialized depending on the use-case (for YAML storage, you'd define `MarshalYAML`).

### Module state

So now, let's define the state of our module.
It is defined in the file `state.go`.
First, change the first line to `package calendar`.
Now, let's go through the code and discover what we need to add:

---

```Go
type state struct {
  Date UniversityDate
}
```

Here, we define that our state holds a UniversityDate.
The `Date` field must be publicly visible so that it can be serialized properly.

---

```Go
type endpoint struct {
	*state
}
```

Keep this intact.
It is the endpoint we will later use to change the date via the web interface.
The endpoint has no internal state apart from a reference to the state it works on.

---

```Go
func newState(input *yaml.Node, ctx api.ServerContext,
	ms api.MessageSender) (api.ModuleState, error) {
  s := new(state)
  if input == nil {
    return s, nil
  }
	if err := input.Decode(&s.Date); err != nil {
    ms.Warning("unable to load UniversityDate: " + err.Error())
    s.Date = UniversityDate(0)
  }
	return s, nil
}
```

Here, we define the constructor for a state object, which is created from YAML input.
YAML is the file format all persistent data is stored in.
A scene's state file contains the data of all visible modules, so we get a part of that file as `yaml.Node` which holds the value for our state.
You do not need to know details about YAML.

`input` may be **`nil`** if the currently stored state has no information about our module.
This will always be the case after adding a module to a scene or loading a new group the first time, so we need to deal with it.
Here, we just return a state with the default value (which will be 0, corresponding to Ick 1st, year 0).

We then try to decode the given YAML node into a date.
An error at this points means that the input data is corrupted.
If that is the case, we issue a warning and load the default value.
Returning an error from the module constructor will halt the main app, so don't do it as long as you can load some default value!

---

```Go

func (s *state) WebView(ctx api.ServerContext) interface{} {
  return s.Date
}

func (s *state) PersistingView(ctx api.ServerContext) interface{} {
	return s.Date
}
```

Now come the serialization functions.
`WebView` returns the data that should be serialized to JSON and send to the web client.
This will use `UniversityDate`'s `MarshalJSON` method.

In `PersistingView`, we need to give the same data we `Decode` the input to in the constructor.
This is the data that will be written to the scene state on the file system.

---

```Go
func (s *state) CreateRendererData() interface{} {
	return s.Date
}
```

This function defines the data we send to the renderer so that it can rebuild its state (e.g. when the group is loaded or the scene changes).
As the renderer runs in another thread, it has its own state and cannot access the `state` object.

We must be careful not to send a pointer into `state` here; that would cause havoc as multiple threads can access the same unprotected data concurrently.
Since `Date` is not a pointer, we are safe here.

---

```Go
func (s *state) PureEndpoint(index int) api.ModulePureEndpoint {
	if index != 0 {
		panic("Endpoint index out of bounds")
	}
	return endpoint{s}
}
```

This function creates our endpoint.
If the module has multiple endpoints, the `index` identifies which endpoint we need to create.
Since we only have one endpoint, we can assume that index is always `0`.

---

```Go
func (e endpoint) Post(payload []byte) (interface{}, interface{},
	api.SendableError) {
  var daysDelta int
  if err := api.ReceiveData(payload, &daysDelta); err != nil {
		return nil, nil, err
	}
  e.state.Date = e.state.Date.add(daysDelta)

	// first value is sent back to client as JSON.
	// second value is sent to Renderer.InitTransition.
	return e.state.Date, e.state.Date, nil
}
```

Finally, this is our endpoint.
We receive the delta (in days) we want to change, and simply apply it to our date.
`api.ReceiveData` is a helper function that wraps any error into an `api.SendableError` and can do some additional validation (which we do not need here).

We send the full date back to the Web Client and also to the Renderer.
This is the same data object we also return in `CreateRendererData`, but that is not required.
Depending on how the animation works, it could be better to send information specifically defining *what* changed (e.g. the day delta).
It could also be that we want to send something different to the client.

### Module Renderer

Let's now implement the module's renderer in the file `renderer.go`.
As before, you need to change the first line to `package calendar`.

---

```Go
type config struct {
	Font       *api.SelectableFont
	Background *api.SelectableTexturedBackground
}
```

This struct defines all editable configuration values of our module.
The struct requires each field's type to be a pointer to an `api.ConfigItem`.
While we could implement our very own config item types here, we'll only be using predefined item types.
If you define custom config item types, you'll also need to add JavaScript to describe how these can be edited in the Web Client.

The type item types we're using are:

 * `SelectableFont`, which allows the user to select the font with which text is rendered.
   It is advicable to use this if you render any text; this enables the user to use any font they like.
 * `SelectableTexturedBackground` allows the user to customize the background color.

Generally, it is advisable to not hard-code colors or fonts into a plugin, so that the user can customize the appearance to match with the configured look of other plugins.

---

```Go
// Renderer implements the rendering of the module's state with SDL.
type Renderer struct {
	config         *config // holds current merged configuration values
  curTex, oldTex *sdl.Texture
  mask           *sdl.Texture
  cur            Date
  oldPos         int32
}
```

Here starts the interesting part:
We need to figure out what data need to be kept as state of the Renderer.
It is advicable to set up the renderer in a way that the `Render` routine is as fast as possible.
This means that we'll pre-render everything into textures and just copy them to the output when `Render` is called.
During animation, we'll need both the old and the new texture, so we put two variables into the `struct` for keeping them.

Additionally, we keep a `mask` texture that contains the current background mask if it is not single-colored.
As you have seen with the HeroList and Title modules, background color selection can include a texture which adds a second color via a given grayscale image (texture) used as alpha channel.
QuestScreen contains functionality to load such an image into a texture with the input image as alpha channel and a separately selected fixed color in the color channels.
Since that texture only depends on the configuration, we keep it as long as the configuration doesn't change.

Finally, `cur` contains the currently displayed date.
We need to remember this in the case that the configuration changes but the data does not to regenerate the currently displayed image.
`oldPos` is our animation state, we'll discuss it when implementing animation.

---

```Go
func newRenderer(backend *sdl.Renderer,
	ms api.MessageSender) (api.ModuleRenderer, error) {
	return &Renderer{}, nil
}

// Descriptor returns the module's metadata
func (*Renderer) Descriptor() *api.Module {
	return &Descriptor
}
```

These are trivial funcs we need to implement.
`newRenderer` initializes the Renderer; we don't need to do anything here since everything will be initialized in `Rebuild`.
Consult the diagram in the [Documentation](/plugins/documentation/#The%20Render%20Loop) for details on the order in which these funcs are called.

---

We'll now create a helper function that creates a texture containing the current date:

```Go
func (r *Renderer) createDateSheet(ctx api.RenderContext,
	d UniversityDate) *sdl.Texture {
	str := fmt.Sprintf("%d %s %d", d.dayOfMonth(), d.month(), d.year())
	face := ctx.Font(
		r.config.Font.FamilyIndex, r.config.Font.Style, r.config.Font.Size)
	strTexture := ctx.TextToTexture(str, face, sdl.Color{R: 0, G: 0, B: 0, A: 255})
	_, _, strWidth, strHeight, _ := strTexture.Query()
	bgColor := r.config.Background.Primary.WithAlpha(255)
	canvas := ctx.CreateCanvas(strWidth+2*ctx.Unit(), strHeight+2*ctx.Unit(),
		&bgColor, r.mask, api.East|api.South|api.West)
	ctx.Renderer().Copy(strTexture, nil, &sdl.Rect{
		X: 2 * ctx.Unit(), Y: ctx.Unit(), W: strWidth, H: strHeight})
	return canvas.Finish()
}
```

Of course, we format the date like a sane person would do: *day month year*.
We query the currently selected font face from the context and render our date string to a texture.
Then, we create a *canvas* based on the dimensions of the rendered text.
A canvas is used to render stuff to a texture we can later use.
We give the inner dimensions, using the text dimension plus one `ctx.Unit()` on each side.
`ctx.Unit()` is a scaling unit that depends on the screen size (just like font sizes do).

`CreateCanvas` optionally renders a background color and mask on it.
We give the selected color and the current mask to it so that it does that for us.
Finally, we select where we want to have borders.
Since we will anchor our date at the top edge of the screen, we create borders for the other three directions.

With the canvas in place, all rendering on the context's `sdl.Renderer` is done into the canvas.
So we copy the rendered text with the renderer, offsetting it so that it is centered (note that we need `2*ctx.Unit()` in x direction since there is a border and another unit of padding).

Finally, we finish the canvas and return the result.

---

```Go
// Rebuild rebuilds the state from the given config and optionally data.
func (r *Renderer) Rebuild(
	ctx api.ExtendedRenderContext, data interface{}, configVal interface{}) {
	r.config = configVal.(*config)
	ctx.UpdateMask(&r.mask, *r.config.Background)
	if data != nil {
		r.cur = data.(UniversityDate)
	}
	if r.curTex != nil {
		r.curTex.Destroy()
	}
	r.curTex = r.createDateSheet(ctx, r.cur)
}
```

This is our `Rebuild` function.
First, we assign the given `config` value.
Then, we update our mask based on it.
we update our data if data is given, and destroy the current texture if it exists.
Finally, we use our helper function to create the new texture.

---

Now comes the animation.
We'll animate date transition like ripping off a sheet from your calendar.
This makes things easy for us:
We only need the old and new date images, and make the old move and fade out while the new fades in.

```Go
// InitTransition starts transitioning after user input changed the state.
func (r *Renderer) InitTransition(
	ctx api.RenderContext, data interface{}) time.Duration {
	r.oldTex = r.curTex
	r.oldTex.SetBlendMode(sdl.BLENDMODE_BLEND)
	r.cur = data.(UniversityDate)
	r.curTex = r.createDateSheet(ctx, r.cur)
	r.curTex.SetBlendMode(sdl.BLENDMODE_BLEND)
	r.oldPos = 0
	return time.Second / 2
}
```

Here, we update internal state from the new data, preserve the current texture in `oldTex` and render the new date into a texture.
We temporarily activate blending for both textures to be able to do fade-in / fade-out.
The initial animation state will be the old texture being completely visible and at the original position.

---

```Go
// TransitionStep advances the transitioning animation.
func (r *Renderer) TransitionStep(
	ctx api.RenderContext, elapsed time.Duration) {
	pos := api.TransitionCurve{Duration: time.Second / 2}.Cubic(elapsed)
	newAlpha := uint8(pos * 255)
	r.oldTex.SetAlphaMod(255 - newAlpha)
	r.curTex.SetAlphaMod(newAlpha)
	_, _, _, oldHeight, _ := r.oldTex.Query()
	r.oldPos = int32(pos * float32(oldHeight) * 3)
}
```

When advancing the animation, we use a `TransitionCurve`, which implements a function going from `0.0` at the beginning to `1.0` at the end of the animation.
The trivial curve would be the linear one, but that looks very boring.
The `Cubic` one we use starts slow, speeds up, and decelerates at the end.
We set the texture's alpha mod to facilitate fading, and the `oldPos` defines how far down the old image is.
We use the image's height for defining how far it moves.

---

```Go
// FinishTransition finalizes the transitioning animation.
func (r *Renderer) FinishTransition(ctx api.RenderContext) {
	r.oldTex.Destroy()
	r.oldTex = nil
	r.curTex.SetBlendMode(sdl.BLENDMODE_NONE)
}
```

At the end of the animation, we destroy the old texture and disable blending for the new one.
We do not need to reset `oldPos` since that is not used outside of animation.

---

```Go
// Render renders the current state / animation frame.
func (r *Renderer) Render(ctx api.RenderContext) {
	sr := ctx.Renderer()
	screenWidth, _, _ := sr.GetOutputSize()
	_, _, curWidth, curHeight, _ := r.curTex.Query()
	sr.Copy(r.curTex, nil, &sdl.Rect{X: screenWidth - curWidth - 5*ctx.Unit(),
		Y: 0, W: curWidth, H: curHeight})
	if r.oldTex != nil {
		_, _, oldWidth, oldHeight, _ := r.oldTex.Query()
		sr.Copy(r.curTex, nil, &sdl.Rect{X: screenWidth - oldWidth - 5*ctx.Unit(),
			Y: r.oldPos, W: oldWidth, H: oldHeight})
	}
}
```

Finally, rendering.
We render the calender to the upper right corner, with a distance of 5 units from the right edge.
If `r.oldTex` is not nil, we're currently animating so we need to render the old date as well.
Fading and position assignment has already been handled by `TransitionStep`.

This wraps up our `renderer.go`.

### Module Descriptor

Now let's go into the `descriptor.go` file and setup or module's descriptor:

```Go
package calendar

import "github.com/QuestScreen/api"

// Descriptor describes this module.
var Descriptor = api.Module{
	Name:                "Calendar",
	ID:                  "tutorial-calendar",
	ResourceCollections: []api.ResourceSelector{},
	EndpointPaths: []string{
		"", // endpoint with no further path; handles updating the date
	},
	DefaultConfig: &config{Font: &api.SelectableFont{
		FamilyIndex: 0, Size: api.HeadingFont, Style: api.Bold},
		Background: &api.SelectableTexturedBackground{
			Primary:      api.RGBColor{Red: 255, Green: 255, Blue: 255},
			TextureIndex: -1,
		},
	},
	CreateRenderer: newRenderer,
	CreateState:    newState,
}
```

The empty endpoint path means that our endpoint is reachable at `/state/tutorial-calendar`.
Since we only have one endpoint, we don't need subpaths.
The `ID` must be unique among all modules, which is why it's usually a good idea to prefix it with the plugin name.
The `Name` is always displayed together with the plugin name, so `"Calendar"` is sufficient here.

The `DefaultConfig` is the config used when the user does not select something else.
As we can't know which fonts the user has installed, we simply use the first one.
We want the text large and bold by default.
For the background, we set the primary color to white and the `TextureIndex` to -1 meaning *no texture*.

### Module UI

The last thing we need to do is to implement the module's UI for changing its state in the Web Client.

 [1]: https://golang.org/
 [2]: https://www.libsdl.org/download-2.0.php
 [3]: https://wiki.lspace.org/mediawiki/Discworld_calendar