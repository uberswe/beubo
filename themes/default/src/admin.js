import '../node_modules/trumbowyg/dist/ui/trumbowyg.css'
import './admin.css'

// Any content field should use WYSIWYG editor
$ ('#contentField').trumbowyg ({
    svgPath: '/default/js/ui/icons.svg'
});