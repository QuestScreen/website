Now that the code has been written, we need to build the plugin.
Since our package has a different name (`discworld` if you followed the tutorial), we need to update our metadata.

| Replace `plugin-tutorial` with `discworld` in the `Makefile`.
|
| Edit the first line of `go.mod` to contain the correct path to your module (ending with `discworld`).

When that is done, you are able to compile your plugin.

| Execute `make` to compile the plugin.

This should produce a file `discworld.so`.
It will be quite large since Go links its entire standard library into it; sadly, linking against a shared standard library with plugins is [currently not supported][1] so there's no way around it.

For debugging, you can also use `make discworld_debug.so` which produces a plugin that can be loaded in a debugging environment.
This works e.g. with VS Code.

To test the plugin, we need QuestScreen installed.
Plugins are placed here:

    ~/.local/share/questscreen/plugins

| If that directory does not exist, create it.
| Place `discworld.so` in that directory.

Now when you start QuestScreen, you sould see the Calendar module in the list of modules.
Let's test it!

| In *Datasets*, create a new group using our *Discworld* group template.
| Then, switch to the group.

You should see **1 Ick 0** in the top right corner.
This is our default date.

| Mess around with the date changing buttons and check the animation.

Notice that month stepping changes the day if it steps over Ick.
Try to fix that if you want!

It seems pretty boring, all white-on-white.

| Go to the group's configuration and set a background color for the calendar.
| Also, back in the group state, set a background image if you want.

Now you can see a colored calendar on top of an image background.

---

This concludes the plugin tutorial.
If you didn't already, I suggest reading the more descriptive and less hands-on [documentation of the plugin API][2] before you start writing your own plugins.
It explains the API's concepts in more depth than we covered here.

 [1]: https://github.com/golang/go/issues/18671
 [2]: /plugins/documentation/