package com.example.quantasonaapp.ui.main

import androidx.activity.ComponentActivity
import androidx.compose.ui.test.junit4.createAndroidComposeRule
import androidx.compose.ui.test.onNodeWithText
import org.junit.Before
import org.junit.Rule
import org.junit.Test

/** UI tests for [com.example.quantasonaapp.ui.main.MainScreen]. */
class MainScreenTest {

  @get:Rule val composeTestRule = createAndroidComposeRule<ComponentActivity>()

  @Before
  fun setup() {
    composeTestRule.setContent { MainScreen(onItemClick = {}) }
  }

  @Test
  fun appHeader_exists() {
    composeTestRule.onNodeWithText("QUANTASONA").assertExists()
  }

  @Test
  fun tabItems_exist() {
    composeTestRule.onNodeWithText("HPA Atlas").assertExists()
    composeTestRule.onNodeWithText("Gem Match").assertExists()
    composeTestRule.onNodeWithText("Geology").assertExists()
    composeTestRule.onNodeWithText("Node HUD").assertExists()
  }
}
