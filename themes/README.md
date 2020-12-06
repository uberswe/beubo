# Themes

This is the themes directory. Beubo will look here for any static files that are needed when loading
pages. I have included third party files directly in the repository and have opted not to use a package
manager like npm because I feel that it's simpler and easier to understand.

The default theme for Beubo can be found here: [https://github.com/uberswe/beubo-default](https://github.com/uberswe/beubo-default)

Any third party libraries and files will have a license file or show a license at the top of the file.

## Best practice

This is a general rule that I like to stick to. Themes should aim to be less than 50kb in size 
per page including css and any js but not counting images. Images should be as small as possible 
(50kb or less) without losing too much quality.