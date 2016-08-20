package Game
import(
	
)
func CalculatePoint(NumArray [13]int,ColorArray [4]int){
	var PlayerPoint [6]int{0,0,0,0,0,0}
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
			
			}else if (NumArray[i] == 2 && PlayerPoint[0] == 4)||(NumArray[i] == 3 && PlayerPoint[0] == 2){
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
			}else if   {
				//同花色閃牌
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