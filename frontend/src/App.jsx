import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import NavBar from './components/Navbar'
import Container from './components/Container'

function App() {
  return (
    <>
      <div id="home-page">
        <NavBar />
        <div id="home-double-container">
          <Container id="box-blue">
            <h1>Premier Container</h1>
          </Container>
          <Container id="box-pink">
            <h1>Deuxième Container</h1>
          </Container>
        </div>
      </div>
    </>
  )
}

export default App
