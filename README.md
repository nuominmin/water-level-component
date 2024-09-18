# water-level-component 水位组件

### 这是一个可以集成到其他项目中作为水位组件



### 通过 GROM 迁移模式同步表结构
```go
    err = db.AutoMigrate(&waterlevel.WaterLevel{})
    if err != nil {
        return nil, nil, err
    }
```


### 示例
```go
    var waterLevelManager waterlevel.Manager
    if waterLevelManager, err = waterlevel.NewManager(db, TASK_NAME); err != nil {
        panic("failed to init water level manager: " + err.Error())
    }
    
    // 历史水位
    var currWaterLevel uint64
    if currWaterLevel, err = waterLevelManager.Load(context.Background()); err != nil {
        panic("failed to load water level: " + err.Error())
    }
    
    if err = waterLevelManager.Save(context.Background(), currWaterLevel); err != nil{
        panic("failed to save water level: " + err.Error())
    }
```
