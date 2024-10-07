import React from 'react'
import {createRoot} from 'react-dom/client'
import App from './App'
import {NextUIProvider} from "@nextui-org/react";
import './style.css'

const container = document.getElementById('root')

const root = createRoot(container!)

root.render(
    <React.StrictMode>
        <NextUIProvider>
            <App/>
        </NextUIProvider>
    </React.StrictMode>
)
