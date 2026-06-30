package com.example.quantasonaapp.domain

import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test

class TesseractGeneratorTest {

    @Test
    fun testTesseractGeneration_lengthAndAlphanumeric() {
        val testData = "test voice anomaly".toByteArray()
        val tesseract = TesseractGenerator.generateTesseractHash(testData)
        
        // Assert exactly 81 characters
        assertEquals(81, tesseract.length)
        
        // Assert all characters are alphanumeric
        assertTrue(tesseract.matches("^[a-zA-Z0-9]+$".toRegex()))
    }
    
    @Test
    fun testTesseractGeneration_deterministic() {
        val testData = "consistent biometric signature".toByteArray()
        val tesseract1 = TesseractGenerator.generateTesseractHash(testData)
        val tesseract2 = TesseractGenerator.generateTesseractHash(testData)
        
        assertEquals(tesseract1, tesseract2)
    }

    @Test
    fun testTesseractGeneration_differentInputs() {
        val testData1 = "signal_alpha".toByteArray()
        val testData2 = "signal_beta".toByteArray()
        
        val tesseract1 = TesseractGenerator.generateTesseractHash(testData1)
        val tesseract2 = TesseractGenerator.generateTesseractHash(testData2)
        
        assertTrue(tesseract1 != tesseract2)
    }
}
