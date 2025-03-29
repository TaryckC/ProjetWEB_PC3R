import { initializeApp } from "firebase/app";
import { getAuth } from "firebase/auth";
import { GoogleAuthProvider } from "firebase/auth";

const firebaseConfig = {
    apiKey: "AIzaSyA4guKkesYflP0dzTPmah4h0W1pR2GpbLA",
    authDomain: "pc3rprojet-ce4a7.firebaseapp.com",
    projectId: "pc3rprojet-ce4a7",
    storageBucket: "pc3rprojet-ce4a7.firebasestorage.app",
    messagingSenderId: "13681784968",
    appId: "1:13681784968:web:1806c86172767f6695c73d",
    measurementId: "G-C1XXMG6VT8"
  };


// Initialize Firebase
const app = initializeApp(firebaseConfig);


// Initialize Firebase Authentication and get a reference to the service
export const auth = getAuth(app);
export const provider = new GoogleAuthProvider();

