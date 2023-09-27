import React from 'react'
import { Global } from '@emotion/react'
import tw, { css, GlobalStyles as BaseStyles } from 'twin.macro'

const customStyles = css`
@import url('https://fonts.googleapis.com/css?family=Comic+Neue:300,400,700|Inter:100,300,500,700,900&display=swap');
  #root {
    ${tw`antialiased `};
    --darkpgred:rgb(130, 10, 0);
    --pgred:rgb(210 70 58);
    --sky:#A3E9FF;
    --darksky:rgb(21 94 117);
    --night:#0F3F4E;
    --darknight:rgb(0 0 0 / 80%);
  }
  html {
    ${tw`bg-transparent`}
  }
  
  h1 {
    ${tw`text-4xl font-bold text-center`}
  }

  body {
    background-color: #333;
    color: white;
  }
  input {
    background-color: #222;
  }
  @font-face {
    font-family: 'Comic Neue';
    src: local('ComicNeue'), url('fonts/ComicNeue-Regular.ttf') format('truetype');
  }
  @font-face {
    font-family: 'Open Sans';
    src: local('Open Sans'), url('fonts/OpenSans-VariableFont.ttf') format('truetype');
  }
  @font-face {
    font-family: 'Silkscreen';
    src:  local('Silkscreen'), url('fonts/Silkscreen-Regular.ttf') format('truetype');
    font-weight: 400;
    font-style: normal;
  }
  @font-face {
    font-family: 'Silkscreen';
    src:  local('Silkscreen'), url('fonts/Silkscreen-Bold.ttf') format('truetype');
    font-weight: 700;
    font-style: normal;
  }
  @font-face {
    font-family: 'PressStart2P';
    src:  local('PressStart2P'), url('fonts/PressStart2P-Regular.ttf') format('truetype');
    font-weight: 400;
    font-style: normal;
  }
  @font-face {
    font-family: 'Pixel Font7';
    src:  local('Pixel Font7'), url('fonts/pixel_font-7.ttf') format('truetype');
    font-weight: 400;
    font-style: normal;
  }
  @font-face {
    font-family: 'munro';
    src:  local('munro'), url('fonts/munro.ttf') format('truetype');
    font-style: normal;
  }
  @font-face {
    font-family: 'Inter';
    src:  local('Inter'), url('fonts/Inter-VariableFont_slnt,wght.ttf') format('truetype');
    font-weight: 100;
  }
  @font-face {
    font-family: 'Inter';
    src:  local('Inter'), url('fonts/Inter-VariableFont_slnt,wght.ttf') format('truetype');
    font-weight: 200;
  }
  @font-face {
    font-family: 'Inter';
    src:  local('Inter'), url('fonts/Inter-VariableFont_slnt,wght.ttf') format('truetype');
    font-weight: 300;
  }
  @font-face {
    font-family: 'Inter';
    src:  local('Inter'), url('fonts/Inter-VariableFont_slnt,wght.ttf') format('truetype');
    font-weight: 400;
  }
  @font-face {
    font-family: 'Inter';
    src:  local('Inter'), url('fonts/Inter-VariableFont_slnt,wght.ttf') format('truetype');
    font-weight: 500;
  }
  @font-face {
    font-family: 'Inter';
    src:  local('Inter'), url('fonts/Inter-VariableFont_slnt,wght.ttf') format('truetype');
    font-weight: 600;
  }
  @font-face {
    font-family: 'Inter';
    src:  local('Inter'), url('fonts/Inter-VariableFont_slnt,wght.ttf') format('truetype');
    font-weight: 700;
  }
  @font-face {
    font-family: 'Inter';
    src:  local('Inter'), url('fonts/Inter-VariableFont_slnt,wght.ttf') format('truetype');
    font-weight: 800;
  }
  @font-face {
    font-family: 'Inter';
    src:  local('Inter'), url('fonts/Inter-VariableFont_slnt,wght.ttf') format('truetype');
    font-weight: 900;
  }
// font-style: normal;

`

const GlobalStyles = () => (
  <>
    <BaseStyles />
    <Global styles={customStyles} />
  </>
)

export default GlobalStyles
