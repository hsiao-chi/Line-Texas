package DB
import (
	_"github.com/go-sql-driver/mysql"
)
//numarray 每個數的數量  numarray[0] 為 2
//colorarray 每個花色的數量 colorarray[0] 為 梅花
func CalculatePoint (ColorArray [4]int,NumArray [13]int) [6]int{
	var PlayerPoint [6]int
	Straight := 0
	for i := 0; i < 4; i++ {
		if ColorArray[i] >=5{
			PlayerPoint[0]=6
		}
	}
	for i := 0; i < 13; i++ {
		if NumArray[i]!=0 {
			Straight++
			if Straight == 5{
				PlayerPoint[0]=5
				PlayerPoint[1]=i
			}else if NumArray[i] == 4 && PlayerPoint[0]<8{
				PlayerPoint[0]=8
				PlayerPoint[5]=PlayerPoint[4]
				PlayerPoint[4]=PlayerPoint[3]
				PlayerPoint[3]=PlayerPoint[2]
				PlayerPoint[2]=PlayerPoint[1]
				PlayerPoint[1]=i
			
			}else if NumArray[i] == 2 && PlayerPoint[0] == 4{
			//葫蘆
				PlayerPoint[0]=7
				PlayerPoint[2]=i
			}else if NumArray[i] == 3 && PlayerPoint[0] == 2{
				PlayerPoint[0]=7
				PlayerPoint[2]=PlayerPoint[1]
				PlayerPoint[1]=i
				
			}else if PlayerPoint[0]==6{
				//同花色閃牌
				PlayerPoint[0]=6
				PlayerPoint[5]=PlayerPoint[4]
				PlayerPoint[4]=PlayerPoint[3]
				PlayerPoint[3]=PlayerPoint[2]
				PlayerPoint[2]=PlayerPoint[1]
				PlayerPoint[1]=i
			}else if NumArray[i] == 3 && PlayerPoint[0]<4{
				PlayerPoint[0]=4
				PlayerPoint[5]=PlayerPoint[4]
				PlayerPoint[4]=PlayerPoint[3]
				PlayerPoint[3]=PlayerPoint[2]
				PlayerPoint[2]=PlayerPoint[1]
				PlayerPoint[1]=i
			}else if NumArray[i] == 2 && PlayerPoint[0] == 2{
				PlayerPoint[0]=3
				PlayerPoint[5]=PlayerPoint[4]
				PlayerPoint[4]=PlayerPoint[3]
				PlayerPoint[3]=PlayerPoint[2]
				PlayerPoint[2]=PlayerPoint[1]
				PlayerPoint[1]=i
			}else if NumArray[i] == 2 && PlayerPoint[0]<2{
				PlayerPoint[0]=2
				PlayerPoint[5]=PlayerPoint[4]
				PlayerPoint[4]=PlayerPoint[3]
				PlayerPoint[3]=PlayerPoint[2]
				PlayerPoint[2]=PlayerPoint[1]
				PlayerPoint[1]=i			
			}else{
				PlayerPoint[0]=1
				PlayerPoint[5]=PlayerPoint[4]
				PlayerPoint[4]=PlayerPoint[3]
				PlayerPoint[3]=PlayerPoint[2]
				PlayerPoint[2]=PlayerPoint[1]
				PlayerPoint[1]=i
			}
		}else{
			Straight=1
		}
	}
	return PlayerPoint
}
func checkwinner(Player1 [6]int,Player2 [6]int,Player3 [6]int,Player4 [6]int,Player5 [6]int,Player6 [6]int,Player7 [6]int,Player8 [6]int,Player9 [6]int,Player10 [6]int) int{
	var players [11][7]int
	for i := 0; i < 6; i++ {
		players[1][i]=Player2[i]	
	}
	for i := 0; i < 6; i++ {
		players[2][i]=Player3[i]	
	}
	for i := 0; i < 6; i++ {
		players[3][i]=Player4[i]	
	}
	for i := 0; i < 6; i++ {
		players[4][i]=Player5[i]	
	}
	for i := 0; i < 6; i++ {
		players[5][i]=Player6[i]	
	}
	for i := 0; i < 6; i++ {
		players[6][i]=Player7[i]	
	}
	for i := 0; i < 6; i++ {
		players[7][i]=Player8[i]	
	}
	for i := 0; i < 6; i++ {
		players[8][i]=Player9[i]	
	}
	for i := 0; i < 6; i++ {
		players[9][i]=Player10[i]	
	}
	for i := 0; i < 6; i++ {
		players[0][i]=Player1[i]	
	}
	var winner int
	var now1 int
	winner = 0
	for j := 0; j < 6; j++ {
		
		now1 = players[0][j]
		winner = 0
		for i := 0; i < 10; i++ {
			if players[i][6]!=-1 {
				if  players[i][j] > now1{
					now1 = players[1][j]
					players[i][6]=j
					winner=i
				}else if players[i][j] == now1{
					players[i][6]=j
				}else{
					players[i][6]=-1
				}
			}
		}
		
	}
	return winner
}
/*
	 NUM[0] NUM[1] NUM[2] NUM[12] NUM[13]
Co[0]	1							1(梅花有幾張)
Co[1]	1			 1				2(菱形有機張)
Co[2]	1							1(愛心有幾張)
Co[3]	1	 1		 1				3(黑桃有幾張)
Co[4]	4	 1	 	 1 (某個數字有幾張)

					主判斷數	 副判斷數 後判斷數
1 	散牌			0		+N 		+K 		+Z		1大
2 	一對			50		+N 		+K 		+Z		2大1大
3 	兩對			100		+N 		+K 		+Z		2大2次大1後大
4 	三條			150		+N 		+K 		+Z		3大1次大1後大
5 	順子			200		+N 		+K 		+Z		1大
6 	同花 		250		+N 		+K 		+Z		由大到小比
7 	葫蘆			300		+N 		+K 		+Z		3大2大
8 	四條			350		+N 		+K 		+Z		4大1大
9 	同花順 		400		+N 		+K 		+Z		1大
10	同花大順		
*/