# -*- coding: utf-8 -*-
"""
Created on Sun Dec 24 12:04:38 2017

@author: Vartotojas
"""

import cv2
import numpy as np
from grabscreen import grab_screen
import glob
import pyautogui
import time

dino = cv2.imread('dino.jpg', 0)
w_dino, h_dino = dino.shape[::-1]

files = glob.glob ('cacti/*.jpg')   
cacti = []       
for file in files:
    temp=cv2.imread(file, 0)  
    cacti.append(temp)
 
second = time.time() 
start_time = time.time()
obstacle_height = 0
should_crouch = False
jumpdist = 175
is_crouching = False
MAX_SPEED_TIME = 118
DINO_WALKING_HEIGHT = 110

while(True):
    leftest = 1000     
    if(time.time() - second >= 1 and time.time() - start_time < MAX_SPEED_TIME):
        jumpdist += 1.7
        second = time.time() 
   # jump_dist = 170  +  (time.time() - start_time) *1.80   
    pts = []
   # scr = grab_screen(region=(75,250, 750, 450))
    scr = grab_screen(region=(75,250, 350, 450))
   # scr = cv2.cvtColor(scr, cv2.COLOR_BGR2RGB)
    scr_gray = cv2.cvtColor(scr, cv2.COLOR_BGR2GRAY)    
    
    res_dino = cv2.matchTemplate(scr_gray, dino, cv2.TM_CCOEFF_NORMED)
    threshold = 0.8
    loc_dino = np.where(res_dino >= threshold)
    dinoX = 0
    dinoH = 0
    
    for pt in zip(*loc_dino[::-1]):        
        cv2.rectangle(scr, pt, (pt[0] + w_dino, pt[1] + h_dino + 9), (50,205,50), 1)
        dinoX = pt[0] + w_dino
        dinoH = pt[1]
        is_crouching = False
    for cactus in cacti: 
        res = cv2.matchTemplate(scr_gray, cactus, cv2.TM_CCOEFF_NORMED)
        w, h = cactus.shape[::-1]
        loc = np.where(res >= 0.8)
        for pt in zip(*loc[::-1]):
            if(pt in pts):
                continue
            cv2.rectangle(scr, pt, (pt[0] + w, pt[1] + h), (0, 0, 255), 1)
            pts.append(pt)
            if(leftest > pt[0] + w and pt[0] > dinoX):
                leftest = pt[0] + w
                obstacle_height = h
                if(h < 20 and pt[1] + 5 < dinoH):
                    should_crouch = True
                elif(should_crouch == True):
                    should_crouch = False
                    pyautogui.keyUp('down')
    
    if(leftest - dinoX < jumpdist - obstacle_height and should_crouch == False):            
        if(dinoH > DINO_WALKING_HEIGHT): 
            cv2.putText(scr, 'Jump!', (0, 20), cv2.FONT_HERSHEY_SIMPLEX, 0.5, (255, 0, 0),
                     2, cv2.LINE_AA)
            pyautogui.press('space')
    elif(should_crouch == True):
        if(is_crouching == False):
           pyautogui.keyDown('down')
           is_crouching = True
        cv2.putText(scr, 'Crouch!', (0, 20), cv2.FONT_HERSHEY_SIMPLEX, 0.5, (30, 255, 10),
                     2, cv2.LINE_AA) 
  
    cv2.imshow('screen', scr) 
    if cv2.waitKey(1) & 0xFF == ord('q'):
        cv2.destroyAllWindows()
        break
    
