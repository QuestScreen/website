Now that the code has been written, we need to build the plugin.
Since our package has a different name (`discworld` if you followed the tutorial), we need to update our metadata.

| Replace all occurences of `PluginTemplate.so` with `discworld.so` in the `Makefile`.
|
| Edit the first line of `go.mod` to contain the correct path to your module (ending with `discworld`).

When that is done, you are able to compile your plugin.

| Execute `make` to compile the plugin.

This should produce a file `discworld.so`.
It will be quite large since Go links its entire standard library into it; sadly, linking against a shared standard library with plugins is [currently not supported][1] so there's no way around it.

To test the plugin, we need QuestScreen installed.
After running QuestScreen once, we'll have the directory structure of the configuration.
We want to place the file `discworld.so` here:

    ~/.local/share/questscreen/plugins


 [1]: https://github.com/golang/go/issues/18671