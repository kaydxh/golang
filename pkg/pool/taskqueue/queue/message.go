/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */

package queue

type Message struct {
	// 业务标识
	Id string

	//内部标识,内部自动生成
	InnerId string
	Name    string

	// the same Scheme for the same task handler
	Scheme string

	// Args is json format, for run task
	Args string
}

type MessageStatus = string

const (
	MessageStatus_Unknown MessageStatus = "Status_Unknown"
	MessageStatus_Success MessageStatus = "Status_Success"
	MessageStatus_Doing   MessageStatus = "Status_Doing"
	MessageStatus_Fail    MessageStatus = "Status_Fail"
)

type MessageResult struct {
	// 业务标识
	Id string

	//内部标识,内部自动生成
	InnerId string
	Name    string

	// the same Scheme for the same task handler
	Scheme string

	// Args is json format, for run task
	Result string

	Status MessageStatus

	Err error
}
