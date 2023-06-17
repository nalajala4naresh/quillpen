call plug#begin('~/.vim/plugged')  " Specify the plugin installation directory

Plug 'fatih/vim-go'                 " Vim-go plugin for Go programming language
Plug 'sheerun/vim-polyglot'
Plug 'tpope/vim-fugitive'           " Vim-fugitive for Git integration
Plug 'vim-airline/vim-airline'      " vim-airline for a status/tabline
Plug 'mtscout6/vim-kubernetes'      " k8s plugin
Plug 'maralla/kubernetes-vim'       " vim k8s plugin
Plug 'kien/ctrlp.vim'               " CtrlP for file searching and navigation
Plug 'preservim/nerdtree'           " NERDTree for file explorer
call plug#end()           " End the plugin section

filetype plugin indent on

" NERDTree configuration
syntax enable
set nocompatible

" Set NERDTree as the default explorer on startup
autocmd VimEnter * NERDTree
set autowrite

let g:go_highlight_fields = 1
let g:go_highlight_functions = 1
let g:go_highlight_function_calls = 1
let g:go_highlight_extra_types = 1
let g:go_highlight_operators = 1

let g:go_fmt_autosave = 1
let g:go_fmt_command = "goimports"

" Status line types/signatures
let g:go_auto_type_info = 1

" Run :GoBuild or :GoTestCompile based on the go file
function! s:build_go_files()
  let l:file = expand('%')
  if l:file =~# '^\f\+_test\.go$'
    call go#test#Test(0, 1)
  elseif l:file =~# '^\f\+\.go$'
    call go#cmd#Build(0)
  endif
endfunction
