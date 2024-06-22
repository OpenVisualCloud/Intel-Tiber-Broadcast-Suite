# Sphinx documentation build configuration file
from __future__ import annotations

import os
import re
import time

import sphinx

project = 'Intel® Tiber™ Broadcast Suite'
copyright = '2024, Intel Corporation'
author = 'Intel Corporation'
release = '0.1.0'

# -- General configuration ---------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#general-configuration

# extensions = [
#     "myst_parser", 'sphinx.ext.graphviz'
# ]
extensions = [
    'myst_parser',
    # 'sphinx.ext.autodoc',
    # 'sphinx.ext.doctest',
    # 'sphinx.ext.todo',
    # 'sphinx.ext.autosummary',
    # 'sphinx.ext.extlinks',
    # 'sphinx.ext.intersphinx',
    # 'sphinx.ext.viewcode',
    # 'sphinx.ext.inheritance_diagram',
    # 'sphinx.ext.coverage',
    'sphinx.ext.graphviz',
    'sphinxcontrib.mermaid'
]
coverage_statistics_to_report = coverage_statistics_to_stdout = True

# use language set by highlight directive if no language is set by role , "sphinx.ext.graphviz"
inline_highlight_respect_highlight = False

# use language set by highlight directive if no role is set
inline_highlight_literals = False

templates_path = ['_templates']
exclude_patterns = ['_build/*', 'tests/*', 'patches/*', 'Thumbs.db', '.DS_Store']


# -- Options for HTML output -------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#options-for-html-output

html_theme = 'sphinx_book_theme'
html_static_path = ['../images']

language = "en"
myst_html_meta = {
    "description lang=en": "Intel® Tiber™ Broadcast Suite",
    "keywords": "Intel®, Intel, Tiber™, Tiber, st20, st22",
    "property=og:locale":  "en_US"
}
myst_fence_as_directive = [ "mermaid" ]

source_suffix = {
    '.rst': 'restructuredtext',
    '.txt': 'restructuredtext',
    '.md': 'markdown',
}
suppress_warnings = ["myst.xref_missing"]

import os
import sys
sys.path.insert(0, os.path.abspath('..'))
sys.path.insert(0, os.path.abspath('../../'))
