import Tagify from '@yaireo/tagify'
import jquery from 'jquery'
import trumbowyg from 'trumbowyg'

import '../node_modules/trumbowyg/dist/ui/trumbowyg.css'
import '../node_modules/@yaireo/tagify/dist/tagify.css'
import './admin.css'

// Any content field should use WYSIWYG editor
$ ('#contentField').trumbowyg ({
    svgPath: '/default/js/ui/icons.svg'
});

let input = document.querySelector ('input[name=tagField]');
if (input !== null) {
    let tagify = new Tagify (input)
}